package websockethub

import(
	"errors"
	"kai-suite/types"
)

var(
	ContactsSyncQueue []types.TxSyncContact
)

func EnqueueContactSync(item types.TxSyncContact) {
	ContactsSyncQueue = append([]types.TxSyncContact{item}, ContactsSyncQueue...)
}

func DequeueContactSync() (item types.TxSyncContact, err error) {
	size := len(ContactsSyncQueue)
	if  size == 0 {
		err = errors.New("Empty")
		return
	}
	item = ContactsSyncQueue[size - 1]
	ContactsSyncQueue = ContactsSyncQueue[:size-1]
	return
}

func GetLastContactSync() (item types.TxSyncContact, err error) {
	size := len(ContactsSyncQueue)
	if size == 0 {
		err = errors.New("Empty")
		return
	}
	item = ContactsSyncQueue[size - 1]
	return 
}
