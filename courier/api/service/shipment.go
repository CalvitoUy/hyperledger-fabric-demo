package service

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/fabricsdk"
	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/inputmodel"
	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/logger"
	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/model"
)

type ShipmentService struct {
	channelClient *channel.Client
	chaincodeID   string
}

func NewShipmentService(channelClient *channel.Client) *ShipmentService {
	chaincodeID := viper.GetString("app.fabric.chaincodeID")

	return &ShipmentService{
		channelClient: channelClient,
		chaincodeID:   chaincodeID,
	}
}

func (this *ShipmentService) UpdateStateAndLocation(trackingCode string, shipmentState *inputmodel.ShipmentState) error {
	organizationPeer0, sellerPeer0, err := fabricsdk.LoadPeers()
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create organization and seller peers")

		return err
	}

	location, _ := json.Marshal(shipmentState.Location)

	updatedAt := time.Now().UTC().Format(time.RFC3339Nano)

	args := [][]byte{
		[]byte(trackingCode),
		[]byte(shipmentState.State),
		location,
		[]byte(strconv.FormatBool(shipmentState.IsDelivered)),
		[]byte(updatedAt),
	}

	response, err := this.channelClient.Execute(
		channel.Request{
			ChaincodeID: this.chaincodeID,
			Fcn:         "changeShipmentStateAndLocation",
			Args:        args,
		},
		channel.WithTargets(organizationPeer0, sellerPeer0), // transaction will be confirmed from both seller and courier1 peers
	)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"response": string(response.Payload),
		}).WithError(err).Error("Unable to update shipment state and location")

		return err
	}

	return nil
}

func (this *ShipmentService) GetAll() ([]model.Shipment, error) {
	args := [][]byte{}
	response, err := this.channelClient.Query(
		channel.Request{
			ChaincodeID: this.chaincodeID,
			Fcn:         "getAllShipments",
			Args:        args,
		},
	)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"response": string(response.Payload),
		}).WithError(err).Error("Unable to get all shipments")

		return nil, err
	}

	records := make(inputmodel.Shipments, 0)
	err = json.Unmarshal(response.Payload, &records)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"payload": string(response.Payload),
		}).WithError(err).Error("Unable to unmarshal all shipments payload")

		return nil, err
	}

	shipments := make([]model.Shipment, 0, len(records))

	for key, value := range records {
		shipment := model.Shipment{}
		err := json.Unmarshal([]byte(value.Value), &shipment)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"key":   key,
				"value": value,
			}).WithError(err).Error("Unable to unmarshal shipment from value")

			return nil, err
		}

		shipments = append(shipments, shipment)
	}

	return shipments, nil
}

func (this *ShipmentService) GetByTrackingCode(trackingCode string) (*model.Shipment, error) {
	args := [][]byte{
		[]byte(trackingCode),
	}
	response, err := this.channelClient.Query(
		channel.Request{
			ChaincodeID: this.chaincodeID,
			Fcn:         "getShipmentByTrackingCode",
			Args:        args,
		},
	)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"response": string(response.Payload),
		}).WithError(err).Errorf("Unable to get shipment %s by tracking code", trackingCode)

		return nil, err
	}

	shipment := model.Shipment{}
	err = json.Unmarshal(response.Payload, &shipment)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"payload": string(response.Payload),
		}).WithError(err).Errorf("Unable to unmarshal shipment %s payload", trackingCode)

		return nil, err
	}

	return &shipment, nil
}

func (this *ShipmentService) GetHistory(trackingCode string) ([]model.Shipment, error) {
	args := [][]byte{
		[]byte(trackingCode),
	}
	response, err := this.channelClient.Query(
		channel.Request{
			ChaincodeID: this.chaincodeID,
			Fcn:         "getShipmentHistory",
			Args:        args,
		},
	)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"response": string(response.Payload),
		}).WithError(err).Errorf("Unable to get shipment %s history", trackingCode)

		return nil, err
	}

	history := make(inputmodel.ShipmentHistory, 0)
	err = json.Unmarshal(response.Payload, &history)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"payload": string(response.Payload),
		}).WithError(err).Errorf("Unable to unmarshal shipment %s history payload", trackingCode)

		return nil, err
	}

	shipments := make([]model.Shipment, 0, len(history))

	for key, value := range history {
		shipment := model.Shipment{}
		err := json.Unmarshal([]byte(value.Value), &shipment)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"key":   key,
				"value": value,
			}).WithError(err).Errorf("Unable to unmarshal shipment % history record from value", trackingCode)

			return nil, err
		}

		shipments = append(shipments, shipment)
	}

	return shipments, nil
}
