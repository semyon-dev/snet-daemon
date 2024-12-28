package training

import (
	"fmt"
	"reflect"

	"github.com/singnet/snet-daemon/v5/storage"
	"github.com/singnet/snet-daemon/v5/utils"
)

type ModelStorage struct {
	delegate storage.TypedAtomicStorage
}

type ModelUserStorage struct {
	delegate storage.TypedAtomicStorage
}

type PendingModelStorage struct {
	delegate storage.TypedAtomicStorage
}

type PublicModelStorage struct {
	delegate storage.TypedAtomicStorage
}

func NewUserModelStorage(atomicStorage storage.AtomicStorage) *ModelUserStorage {
	prefixedStorage := storage.NewPrefixedAtomicStorage(atomicStorage, "/model-user/userModelStorage")
	userModelStorage := storage.NewTypedAtomicStorageImpl(
		prefixedStorage, serializeModelUserKey, reflect.TypeOf(ModelUserKey{}), utils.Serialize, utils.Deserialize,
		reflect.TypeOf(ModelUserData{}),
	)
	return &ModelUserStorage{delegate: userModelStorage}
}

func NewModelStorage(atomicStorage storage.AtomicStorage) *ModelStorage {
	prefixedStorage := storage.NewPrefixedAtomicStorage(atomicStorage, "/model-user/modelStorage")
	modelStorage := storage.NewTypedAtomicStorageImpl(
		prefixedStorage, serializeModelKey, reflect.TypeOf(ModelKey{}), utils.Serialize, utils.Deserialize,
		reflect.TypeOf(ModelData{}),
	)
	return &ModelStorage{delegate: modelStorage}
}

func NewPendingModelStorage(atomicStorage storage.AtomicStorage) *PendingModelStorage {
	prefixedStorage := storage.NewPrefixedAtomicStorage(atomicStorage, "/model-user/pendingModelStorage")
	pendingModelStorage := storage.NewTypedAtomicStorageImpl(
		prefixedStorage, serializePendingModelKey, reflect.TypeOf(PendingModelKey{}), utils.Serialize, utils.Deserialize,
		reflect.TypeOf(PendingModelData{}),
	)
	return &PendingModelStorage{delegate: pendingModelStorage}
}

func NewPublicModelStorage(atomicStorage storage.AtomicStorage) *PublicModelStorage {
	prefixedStorage := storage.NewPrefixedAtomicStorage(atomicStorage, "/model-user/publicModelStorage")
	publicModelStorage := storage.NewTypedAtomicStorageImpl(
		prefixedStorage, serializePublicModelKey, reflect.TypeOf(PublicModelKey{}), utils.Serialize, utils.Deserialize,
		reflect.TypeOf(PublicModelData{}),
	)
	return &PublicModelStorage{delegate: publicModelStorage}
}

type ModelKey struct {
	OrganizationId string
	ServiceId      string
	GroupId        string
	//GRPCMethodName  string
	//GRPCServiceName string
	ModelId string
}

func (key *ModelKey) String() string {
	return fmt.Sprintf("{ID:%v|%v|%v|%v}", key.OrganizationId,
		key.ServiceId, key.GroupId, key.ModelId)
}

type ModelData struct {
	IsPublic            bool
	ModelName           string
	AuthorizedAddresses []string
	Status              Status
	CreatedByAddress    string
	ModelId             string
	UpdatedByAddress    string
	GroupId             string
	OrganizationId      string
	ServiceId           string
	GRPCMethodName      string
	GRPCServiceName     string
	Description         string
	IsDefault           bool
	TrainingLink        string
	UpdatedDate         string
}

func (data *ModelData) String() string {
	return fmt.Sprintf("{DATA:%v|%v|%v|%v|%v|%v|IsPublic:%v|accesibleAddress:%v|createdBy:%v|updatedBy:%v|status:%v|TrainingLin:%v}",
		data.OrganizationId,
		data.ServiceId, data.GroupId, data.GRPCServiceName, data.GRPCMethodName, data.ModelId, data.AuthorizedAddresses, data.IsPublic,
		data.CreatedByAddress, data.UpdatedByAddress, data.Status, data.TrainingLink)
}

type ModelUserKey struct {
	OrganizationId string
	ServiceId      string
	GroupId        string
	//GRPCMethodName  string
	//GRPCServiceName string
	UserAddress string
}

func (key *ModelUserKey) String() string {
	return fmt.Sprintf("{ID:%v|%v|%v|%v}", key.OrganizationId,
		key.ServiceId, key.GroupId, key.UserAddress)
}

// ModelUserData maintain the list of all modelIds for a given user address
type ModelUserData struct {
	ModelIds []string
	//the below are only for display purposes
	OrganizationId string
	ServiceId      string
	GroupId        string
	//GRPCMethodName  string
	//GRPCServiceName string
	UserAddress string
}

func (data *ModelUserData) String() string {
	return fmt.Sprintf("{DATA:%v|%v|%v|%v|%v}",
		data.OrganizationId,
		data.ServiceId, data.GroupId, data.UserAddress, data.ModelIds)
}

type PendingModelKey struct {
	OrganizationId string
	ServiceId      string
	GroupId        string
}

func (key *PendingModelKey) String() string {
	return fmt.Sprintf("{ID:%v|%v|%v}", key.OrganizationId, key.ServiceId, key.GroupId)
}

type PendingModelData struct {
	ModelIDs []string
}

// PendingModelData maintain the list of all modelIds that have TRAINING\VALIDATING status
func (data *PendingModelData) String() string {
	return fmt.Sprintf("{DATA:%v}", data.ModelIDs)
}

type PublicModelKey struct {
	OrganizationId string
	ServiceId      string
	GroupId        string
}

func (key *PublicModelKey) String() string {
	return fmt.Sprintf("{ID:%v|%v|%v}", key.OrganizationId, key.ServiceId, key.GroupId)
}

type PublicModelData struct {
	ModelIDs []string
}

func (data *PublicModelData) String() string {
	return fmt.Sprintf("{DATA:%v}", data.ModelIDs)
}

func serializeModelKey(key any) (serialized string, err error) {
	modelKey := key.(*ModelKey)
	return modelKey.String(), nil
}

func (storage *ModelStorage) Get(key *ModelKey) (state *ModelData, ok bool, err error) {
	value, ok, err := storage.delegate.Get(key)
	if err != nil || !ok {
		return nil, ok, err
	}
	return value.(*ModelData), ok, err
}

func (storage *ModelStorage) GetAll() (states []*ModelData, err error) {
	values, err := storage.delegate.GetAll()
	if err != nil {
		return
	}

	return values.([]*ModelData), nil
}

func (storage *ModelStorage) Put(key *ModelKey, state *ModelData) (err error) {
	return storage.delegate.Put(key, state)
}

func (storage *ModelStorage) PutIfAbsent(key *ModelKey, state *ModelData) (ok bool, err error) {
	return storage.delegate.PutIfAbsent(key, state)
}

func (storage *ModelStorage) CompareAndSwap(key *ModelKey, prevState *ModelData,
	newState *ModelData) (ok bool, err error) {
	return storage.delegate.CompareAndSwap(key, prevState, newState)
}

func serializeModelUserKey(key any) (serialized string, err error) {
	modelUserKey := key.(*ModelUserKey)
	return modelUserKey.String(), nil
}

func (storage *ModelUserStorage) Get(key *ModelUserKey) (state *ModelUserData, ok bool, err error) {
	value, ok, err := storage.delegate.Get(key)
	if err != nil || !ok {
		return nil, ok, err
	}
	return value.(*ModelUserData), ok, err
}

func (storage *ModelUserStorage) GetAll() (states []*ModelUserData, err error) {
	values, err := storage.delegate.GetAll()
	if err != nil {
		return
	}

	return values.([]*ModelUserData), nil
}

func (storage *ModelUserStorage) Put(key *ModelUserKey, state *ModelUserData) (err error) {
	return storage.delegate.Put(key, state)
}

func (storage *ModelUserStorage) PutIfAbsent(key *ModelUserKey, state *ModelUserData) (ok bool, err error) {
	return storage.delegate.PutIfAbsent(key, state)
}

func (storage *ModelUserStorage) CompareAndSwap(key *ModelUserKey, prevState *ModelUserData,
	newState *ModelUserData) (ok bool, err error) {
	return storage.delegate.CompareAndSwap(key, prevState, newState)
}

func serializePendingModelKey(key any) (serialized string, err error) {
	pendingModelKey := key.(*PendingModelKey)
	return pendingModelKey.String(), nil
}

func (storage *PendingModelStorage) Get(key *PendingModelKey) (state *PendingModelData, ok bool, err error) {
	value, ok, err := storage.delegate.Get(key)
	if err != nil || !ok {
		return nil, ok, err
	}

	return value.(*PendingModelData), ok, err
}

func (storage *PendingModelStorage) GetAll() (states []*PendingModelData, err error) {
	values, err := storage.delegate.GetAll()
	if err != nil {
		return
	}

	return values.([]*PendingModelData), nil
}

func (storage *PendingModelStorage) Put(key *PendingModelKey, state *PendingModelData) (err error) {
	return storage.delegate.Put(key, state)
}

func (pendingStorage *PendingModelStorage) AddPendingModelId(key *PendingModelKey, modelId string) (err error) {
	typedUpdateFunc := func(conditionValues []storage.TypedKeyValueData) (update []storage.TypedKeyValueData, ok bool, err error) {
		if len(conditionValues) != 1 || conditionValues[0].Key != key {
			return nil, false, fmt.Errorf("unexpected condition values or missing key")
		}

		// Fetch the current list of pending model IDs from the storage
		currentValue, ok, err := pendingStorage.delegate.Get(key)
		if err != nil {
			return nil, false, err
		}

		var pendingModelData *PendingModelData
		if currentValue == nil {
			pendingModelData = &PendingModelData{ModelIDs: make([]string, 0, 100)}
		} else {
			pendingModelData = currentValue.(*PendingModelData)
		}

		// Check if the modelId already exists
		for _, currentModelId := range pendingModelData.ModelIDs {
			if currentModelId == modelId {
				// If the model ID already exists, no update is needed
				return nil, false, nil
			}
		}

		// Add the new model ID to the list
		pendingModelData.ModelIDs = append(pendingModelData.ModelIDs, modelId)

		// Prepare the updated values for the transaction
		newValues := []storage.TypedKeyValueData{
			{
				Key:     key,
				Value:   pendingModelData,
				Present: true,
			},
		}

		return newValues, true, nil
	}

	request := storage.TypedCASRequest{
		ConditionKeys:           []any{key},
		RetryTillSuccessOrError: true,
		Update:                  typedUpdateFunc,
	}

	// Execute the transaction
	ok, err := pendingStorage.delegate.ExecuteTransaction(request)
	if err != nil {
		return fmt.Errorf("transaction execution failed: %w", err)
	}
	if !ok {
		return fmt.Errorf("transaction was not successful")
	}

	return nil
}

func (storage *PendingModelStorage) PutIfAbsent(key *PendingModelKey, state *PendingModelData) (ok bool, err error) {
	return storage.delegate.PutIfAbsent(key, state)
}

func (storage *PendingModelStorage) CompareAndSwap(key *PendingModelKey, prevState *PendingModelData,
	newState *PendingModelData) (ok bool, err error) {
	return storage.delegate.CompareAndSwap(key, prevState, newState)
}

func serializePublicModelKey(key any) (serialized string, err error) {
	pendingModelKey := key.(*PublicModelKey)
	return pendingModelKey.String(), nil
}

func (storage *PublicModelStorage) Get(key *PublicModelKey) (state *PublicModelData, ok bool, err error) {
	value, ok, err := storage.delegate.Get(key)
	if err != nil || !ok {
		return nil, ok, err
	}

	return value.(*PublicModelData), ok, err
}

func (storage *PublicModelStorage) GetAll() (states []*PublicModelData, err error) {
	values, err := storage.delegate.GetAll()
	if err != nil {
		return
	}

	return values.([]*PublicModelData), nil
}

func (storage *PublicModelStorage) Put(key *PublicModelKey, state *PublicModelData) (err error) {
	return storage.delegate.Put(key, state)
}

func (publicStorage *PublicModelStorage) AddPublicModelId(key *PublicModelKey, modelId string) (err error) {
	typedUpdateFunc := func(conditionValues []storage.TypedKeyValueData) (update []storage.TypedKeyValueData, ok bool, err error) {
		if len(conditionValues) != 1 || conditionValues[0].Key != key {
			return nil, false, fmt.Errorf("unexpected condition values or missing key")
		}

		// Fetch the current list of public model IDs from the storage
		currentValue, ok, err := publicStorage.delegate.Get(key)
		if err != nil {
			return nil, false, err
		}

		var publicModelData *PublicModelData
		if currentValue == nil {
			publicModelData = &PublicModelData{ModelIDs: make([]string, 0, 100)}
		} else {
			publicModelData = currentValue.(*PublicModelData)
		}

		// Check if the modelId already exists
		for _, currentModelId := range publicModelData.ModelIDs {
			if currentModelId == modelId {
				// If the model ID already exists, no update is needed
				return nil, false, nil
			}
		}

		// Add the new model ID to the list
		publicModelData.ModelIDs = append(publicModelData.ModelIDs, modelId)

		// Prepare the updated values for the transaction
		newValues := []storage.TypedKeyValueData{
			{
				Key:     key,
				Value:   publicModelData,
				Present: true,
			},
		}

		return newValues, true, nil
	}

	request := storage.TypedCASRequest{
		ConditionKeys:           []any{key},
		RetryTillSuccessOrError: true,
		Update:                  typedUpdateFunc,
	}

	// Execute the transaction
	ok, err := publicStorage.delegate.ExecuteTransaction(request)
	if err != nil {
		return fmt.Errorf("transaction execution failed: %w", err)
	}
	if !ok {
		return fmt.Errorf("transaction was not successful")
	}

	return nil
}

func (storage *PublicModelStorage) PutIfAbsent(key *PublicModelKey, state *PublicModelData) (ok bool, err error) {
	return storage.delegate.PutIfAbsent(key, state)
}

func (storage *PublicModelStorage) CompareAndSwap(key *PublicModelKey, prevState *PublicModelData,
	newState *PublicModelData) (ok bool, err error) {
	return storage.delegate.CompareAndSwap(key, prevState, newState)
}
