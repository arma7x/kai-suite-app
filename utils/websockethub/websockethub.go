package websockethub

import (
	"strings"
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
	custom_widget "kai-suite/widgets"
)

var (
	address string
	connectionChan = make(chan bool)
	clientConnectedChan = make(chan bool)
	syncProgressChan = make(chan bool)
	progressDialog *custom_widget.ProgressInfiniteDialog
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
	reloadThreadsCb func(map[int]*types.MozMobileMessageThread)
	reloadMessages func(map[int][]*types.MozSmsMessage)
	removeContactCb func(string, *people.Person)
	refreshThreadsCb func()
)

func init() {
	go func () {
		for {
			select {
				case sync := <- syncProgressChan:
					log.Info("syncProgressChan: ", sync)
					if sync == true {
						progressDialog = custom_widget.NewProgressInfinite("Synchronizing", "Please wait...", global.WINDOW)
					} else {
						progressDialog.Hide()
					}
			}
		}
	}()
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") != "websocket" && r.Header.Get("Connection") != "Upgrade" {
		fmt.Fprintf(w, "PC Suite for KaiOS device")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn("upgrade:", err)
		return
	}
	Client = types.CreateClient("Unknown", conn)
	clientConnectedChan <- true
	log.Info("upgrade success")
	defer Client.GetConn().Close()
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			Client = nil
			GoogleContactsQueue = nil
			clientConnectedChan <- false
			log.Warn(err.Error())
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
							data := types.RxSyncDevice0{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								Client.SetDevice(data.Device)
								Client.SetIMEI(data.IMEI)
								log.Info("IMEI: ", Client.GetIMEI())
							}
							clientConnectedChan <- true
						case 2:
							data := types.RxSyncContactFlag2{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								// log.Info(rx.Flag, ": ", data.Namespace, ": ", data.SyncID, ": ", data.SyncUpdated)
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
													log.Error(err.Error())
													return err
												}
												return nil
											}
											log.Error(err.Error())
											return err
										} else {
											log.Error(err.Error())
											return err
										}
										return nil
									})
								}
								SyncGoogleContact()
							}
						case 4:
							data := types.RxSyncContactFlag4{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								// b, _:= json.Marshal(data)
								// log.Info(rx.Flag, ": ", data.Namespace, ": ", string(b), ": ", data.KaiContact.Updated)
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
													log.Error(err.Error())
													return err
												}
												EnqueueContactSync(types.TxSyncGoogleContact{Namespace: data.Namespace, Metadata: metadata, Person: &person}, true)
												return SyncGoogleContact()
											} else {
												log.Error(err.Error())
												return err
											}
											return nil
										}
										log.Error(err.Error())
										return err
									} else {
										log.Error(err.Error())
										return err
									}
									return nil
								}); err != nil {
									SyncGoogleContact()
								}
							}
						case 6:
							data := types.RxSyncContactFlag6{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								split := strings.Split(data.Namespace, ":")
								if len(split) == 3 {
									// log.Info(rx.Flag, ": Delete ", split[0], ":", split[1], ":", split[2])
									global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
										val, err := tx.Get(data.Namespace)
										if err != nil {
											return err
										}
										var person people.Person
										if err := json.Unmarshal([]byte(val), &person); err != nil {
											return err
										}
										if _, err = tx.Delete(data.Namespace); err != nil {
											return err
										}
										metadata := types.Metadata{}
										if metadata_s, err := tx.Get("metadata:" + data.Namespace); err == nil {
											if err := json.Unmarshal([]byte(metadata_s), &metadata); err != nil {
												return err
											}
											metadata.Deleted = true
											if metadata_b, err := json.Marshal(metadata); err == nil {
												tx.Set("metadata:" + data.Namespace, string(metadata_b[:]), nil)
												removeContactCb(split[0], &person)
											} else {
												return err
											}
										} else {
											return err
										}
										return nil
									})
								}
							}
							SyncGoogleContact()
						case 8:
							data := types.RxRestoreContactFlag8{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
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
													log.Error(err.Error())
													return err
												}
												return nil
											}
											log.Error(err.Error())
											return err
										} else {
											log.Error(err.Error())
											return err
										}
										return nil
									})
								}
								RestoreGoogleContact()
							}
						case 10:
							data := types.RxSyncLocalContactFlag10{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								log.Info("PushList: ", len(data.PushList))
								for _, item := range data.PushList {
									global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
										key := "local:people:" + item.KaiContact.Key[0]
										metadataKey := "metadata:local:people:" + item.KaiContact.Key[0]
										person := people.Person{}
										name := &people.Name{}
										person.Names = make([]*people.Name, 1)
										phoneNumber := &people.PhoneNumber{}
										person.PhoneNumbers = make([]*people.PhoneNumber, 1)
										emailAddress := &people.EmailAddress{}
										person.EmailAddresses = make([]*people.EmailAddress, 1)
										if len(item.KaiContact.Name) > 0 {
											name.UnstructuredName = item.KaiContact.Name[0]
										}
										if len(item.KaiContact.GivenName) > 0 {
											name.GivenName = item.KaiContact.GivenName[0]
										}
										if len(item.KaiContact.FamilyName) > 0 {
											name.FamilyName = item.KaiContact.FamilyName[0]
										}
										if len(item.KaiContact.Tel) > 0 {
											if len(item.KaiContact.Tel[0].Type) > 0 {
												phoneNumber.Type = item.KaiContact.Tel[0].Type[0]
											}
											if len(item.KaiContact.Tel[0].Value) > 0 {
												phoneNumber.Value = item.KaiContact.Tel[0].Value
											}
										}
										if len(item.KaiContact.Email) > 0 {
											if len(item.KaiContact.Email[0].Type) > 0 {
												emailAddress.Type = item.KaiContact.Email[0].Type[0]
											}
											if len(item.KaiContact.Email[0].Value) > 0 {
												emailAddress.Value = item.KaiContact.Email[0].Value
											}
										}
										person.Names[0] = name
										person.PhoneNumbers[0] = phoneNumber
										person.EmailAddresses[0] = emailAddress
										person.ResourceName = "people/" + item.KaiContact.Key[0]
										b, _ := person.MarshalJSON()
										hash := sha256.Sum256(b)
										item.Metadata.Hash = hex.EncodeToString(hash[:])
										item.Metadata.Deleted = false
										//log.Info(string(b))
										mb, _ := json.Marshal(item.Metadata)
										//log.Info(string(mb))
										tx.Set(key, string(b), nil)
										tx.Set(metadataKey, string(mb), nil)
										return nil
									})
								}

								log.Info("DeleteList: ", len(data.DeleteList))
								for _, item := range data.DeleteList {
									global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
										key := "local:people:" + item.SyncID
										metadataKey := "metadata:local:people:" + item.SyncID
										tx.Delete(key)
										tx.Delete(metadataKey)
										return nil
									})
								}

								log.Info("SyncList: ", len(data.SyncList))
								for _, item := range data.SyncList {
									key := "local:people:" + item.KaiContact.Key[0]
									metadataKey := "metadata:local:people:" + item.Metadata.SyncID
									// log.Info(key, " : ", metadataKey)
									global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
										val, err := tx.Get(key)
										if err != nil {
											return err
										}
										var person people.Person
										if err := json.Unmarshal([]byte(val), &person); err != nil {
											return err
										}
										if len(item.KaiContact.Name) > 0 {
											person.Names[0].UnstructuredName = item.KaiContact.Name[0]
										}
										if len(item.KaiContact.GivenName) > 0 {
											person.Names[0].GivenName = item.KaiContact.GivenName[0]
										}
										if len(item.KaiContact.FamilyName) > 0 {
											person.Names[0].FamilyName = item.KaiContact.FamilyName[0]
										}
										if len(item.KaiContact.Tel) > 0 {
											if len(item.KaiContact.Tel[0].Type) > 0 {
												person.PhoneNumbers[0].Type = item.KaiContact.Tel[0].Type[0]
											}
											if len(item.KaiContact.Tel[0].Value) > 0 {
												person.PhoneNumbers[0].Value = item.KaiContact.Tel[0].Value
											}
										}
										if len(item.KaiContact.Email) > 0 {
											if len(item.KaiContact.Email[0].Type) > 0 {
												person.EmailAddresses[0].Type = item.KaiContact.Email[0].Type[0]
											}
											if len(item.KaiContact.Email[0].Value) > 0 {
												person.EmailAddresses[0].Value = item.KaiContact.Email[0].Value
											}
										}
										b, _ := person.MarshalJSON()
										hash := sha256.Sum256(b)
										metadata := types.Metadata{}
										if metadata_s, err := tx.Get(metadataKey); err == nil {
											if err := json.Unmarshal([]byte(metadata_s), &metadata); err == nil {
												metadata.SyncID = item.KaiContact.Key[0]
												metadata.SyncUpdated = item.KaiContact.Updated
												metadata.Hash = hex.EncodeToString(hash[:])
												if metadata_b, err := json.Marshal(metadata); err == nil {
													// log.Info(string(metadata_b[:]))
													tx.Set(metadataKey, string(metadata_b[:]), nil)
													b2, _ := person.MarshalJSON()
													// log.Info(string(b2))
													if _, _, err := tx.Set(key, string(b2), nil); err != nil {
														log.Error(err.Error())
														return err
													}
													return nil
												} else {
													log.Error(err.Error())
													return err
												}
												return nil
											}
											log.Error(err.Error())
											return err
										} else {
											log.Error(err.Error())
											return err
										}
										return nil
									})
								}
								log.Info("MergedList: ", len(data.MergedList)) // TODO
								syncProgressChan <- false
							}
						case 12:
							data := types.RxSyncSMSFlag12{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								reloadThreadsCb(data.Threads)
								reloadMessages(data.Messages)
								refreshThreadsCb()
							}
						case 14:
							syncProgressChan <- false
							data := types.RxSyncEvents14{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								StartSyncEvent(data.Namespace, data.UnsyncEvents);
							}
					}
				}
		}
	}
}

func GetConnectionChan() chan bool {
	return connectionChan
}

func GetClientConnectedChan() chan bool {
	return clientConnectedChan
}

func RegisterListener(_reloadThreadsCb func(map[int]*types.MozMobileMessageThread), _reloadMessages func(map[int][]*types.MozSmsMessage), _removeContactCb func(string, *people.Person), _refreshThreadsCb func()) {
	reloadThreadsCb = _reloadThreadsCb
	reloadMessages = _reloadMessages
	removeContactCb = _removeContactCb
	refreshThreadsCb = _refreshThreadsCb
}

func Start(addr string) error {
	address = addr
	Status = true
	connectionChan <- true
	m := http.NewServeMux()
	Server = http.Server{Addr: address, Handler: m}
	m.HandleFunc("/", websocketHandler)
	m.HandleFunc("/local-contacts", localContactListHandler)
	if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Status = false
		connectionChan <- false
		return err
	}
	return nil
}

func Stop() error {
	if err := Server.Shutdown(context.Background()); err != nil {
		return err
	} else {
		if Client != nil {
			GoogleContactsQueue = nil
			Client.GetConn().WriteMessage(websocket.CloseMessage, []byte{})
			Client.GetConn().Close()
			Client = nil
		}
		Status = false
		connectionChan <- false
		clientConnectedChan <- false
	}
	return nil
}
