package websockethub

import (
	"time"
	"fmt"
	"io"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/gorilla/websocket"
	"context"
	"kai-suite/types"
	"encoding/json"
	"kai-suite/utils/global"
	"github.com/tidwall/buntdb"
	"google.golang.org/api/people/v1"
	"crypto/sha256"
	"encoding/hex"
)

var (
	initialized = false
	address string
	clientVisibilityChan chan bool
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	Status bool = false
	Server http.Server
	Client *types.Client
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") != "websocket" && r.Header.Get("Connection") != "Upgrade" {
		fmt.Fprintf(w, "PC Suite for KaiOS device")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	// id as time
	Client = types.CreateClient("Unknown", conn)
	clientVisibilityChan <- true
	log.Info("upgrade success")
	defer Client.GetConn().Close()
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			Client = nil
			ContactsSyncQueue = nil
			clientVisibilityChan <- false
			log.Error(err)
			if websocket.IsCloseError(err, websocket.CloseGoingAway) || err == io.EOF {
				break
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
			break
		}
		switch mt {
			case websocket.TextMessage:
				rx := types.WebsocketMessageFlag{}
				if err := json.Unmarshal(msg, &rx); err == nil {
					switch rx.Flag {
						case 0:
							Client.SetDevice(rx.Data)
							clientVisibilityChan <- true
						case 2:
							data := types.RxSyncContactFlag2{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								log.Info(rx.Flag, ": ", data.Namespace, ": ", data.SyncID, ": ", data.SyncUpdated)
								if data.SyncID != "error" {
									global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
										metadata := types.Metadata{}
										if metadata_s, err := tx.Get("metadata:" + data.Namespace); err == nil {
											if err := json.Unmarshal([]byte(metadata_s), &metadata); err == nil {
												metadata.SyncID = data.SyncID
												metadata.SyncUpdated = data.SyncUpdated
												if metadata_b, err := json.Marshal(metadata); err == nil {
													tx.Set("metadata:" + data.Namespace, string(metadata_b[:]), nil)
												} else {
													log.Warn(err.Error())
													return err
												}
												return nil
											}
											log.Warn(err.Error())
											return err
										} else {
											log.Warn(err.Error())
											return err
										}
										return nil
									})
								}
								if item, err := DequeueContactSync(); err == nil && Client != nil {
									bd, _ := json.Marshal(item)
									btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 1, Data: string(bd)})
									if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
										log.Error("write:", err)
									}
								}
							}
						case 4:
							data := types.RxSyncContactFlag4{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								b, _:= json.Marshal(data)
								log.Info(rx.Flag, ": ", data.Namespace, ": ", string(b), ": ", data.KaiContact.Updated)
								if err := global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
									val, err := tx.Get(data.Namespace)
									if err != nil {
										return err
									}
									var person people.Person
									if err := json.Unmarshal([]byte(val), &person); err != nil {
										return err
									}
									if len(data.KaiContact.Name) > 0 {
										person.Names[0].UnstructuredName = data.KaiContact.Name[0]
									}
									if len(data.KaiContact.GivenName) > 0 {
										person.Names[0].GivenName = data.KaiContact.GivenName[0]
									}
									if len(data.KaiContact.FamilyName) > 0 {
										person.Names[0].FamilyName = data.KaiContact.FamilyName[0]
									}
									if len(data.KaiContact.Tel) > 0 {
										if len(data.KaiContact.Tel[0].Type) > 0 { 
											person.PhoneNumbers[0].Type = data.KaiContact.Tel[0].Type[0]
										}
										if len(data.KaiContact.Tel[0].Value) > 0 { 
											person.PhoneNumbers[0].Value = data.KaiContact.Tel[0].Value
										}
									}
									if len(data.KaiContact.Email) > 0 {
										if len(data.KaiContact.Email[0].Type) > 0 { 
											person.EmailAddresses[0].Type = data.KaiContact.Email[0].Type[0]
										}
										if len(data.KaiContact.Email[0].Value) > 0 { 
											person.EmailAddresses[0].Value = data.KaiContact.Email[0].Value
										}
									}
									person.Metadata.Sources[0].UpdateTime = ""
									b, _ := person.MarshalJSON()
									hash := sha256.Sum256(b)
									metadata := types.Metadata{}
									if metadata_s, err := tx.Get("metadata:" + data.Namespace); err == nil {
										if err := json.Unmarshal([]byte(metadata_s), &metadata); err == nil {
											metadata.SyncID = data.KaiContact.Id
											metadata.SyncUpdated = data.KaiContact.Updated
											metadata.Hash = hex.EncodeToString(hash[:])
											if metadata_b, err := json.Marshal(metadata); err == nil {
												tx.Set("metadata:" + data.Namespace, string(metadata_b[:]), nil)
												person.Metadata.Sources[0].UpdateTime = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
												b2, _ := person.MarshalJSON()
												if _, _, err := tx.Set(data.Namespace, string(b2), nil); err != nil {
													log.Warn(err.Error())
													return err
												}
												EnqueueContactSync(types.TxSyncContact{Namespace: data.Namespace, Metadata: metadata, Person: &person}, true)
												if item, err := DequeueContactSync(); err == nil && Client != nil {
													bd, _ := json.Marshal(item)
													btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 1, Data: string(bd)})
													if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
														log.Warn(err.Error())
														return err
													}
												}
											} else {
												log.Warn(err.Error())
												return err
											}
											return nil
										}
										log.Warn(err.Error())
										return err
									} else {
										log.Warn(err.Error())
										return err
									}
									return nil
								}); err != nil {
									if item, err := DequeueContactSync(); err == nil && Client != nil {
										bd, _ := json.Marshal(item)
										btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 1, Data: string(bd)})
										if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
											log.Error("write:", err)
										}
									}
								}
							}
						case 6:
							log.Info(rx.Flag, ": ", rx.Data)
					}
				}
		}
	}
}

func Init(addr string, clientChan chan bool) {
	initialized = true
	address = addr
	clientVisibilityChan = clientChan
}

func Start(fn func(bool, error)) {
	if initialized == false {
		return
	}
	Status = true
	fn(Status, nil)
	m := http.NewServeMux()
	Server = http.Server{Addr: address, Handler: m}
	m.HandleFunc("/", handler)
	if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Status = false
		fn(Status, err)
	}
}

func Stop(fn func(bool, error)) {
	if initialized == false {
		return
	}
	if err := Server.Shutdown(context.Background()); err != nil {
		fn(Status, err)
	} else {
		if Client != nil {
			ContactsSyncQueue = nil
			Client.GetConn().WriteMessage(websocket.CloseMessage, []byte{})
			Client.GetConn().Close()
			Client = nil
		}
		Status = false
		fn(Status, nil)
	}
}
