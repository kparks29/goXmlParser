package mismo

import (
	"github.com/homdna/homdna-models"
)

func MergeHomdnas(originalHomdna models.HomdnaModel, parsedHomdna models.HomdnaModel) models.HomdnaModel {
	homdna := models.HomdnaModel{}
	homdna = originalHomdna

	// Merge Addresses
	homdna.Address = MergeAddresses(originalHomdna.Address, parsedHomdna.Address)

	// Merge Lot
	homdna.Lot = MergeLot(originalHomdna.Lot, parsedHomdna.Lot)
	// Merge Structures
	homdna.Structures = MergeStructures(originalHomdna.Structures, parsedHomdna.Structures)

	return homdna
}

func MergeAddresses(originalAddress models.AddressModel, parsedAddress models.AddressModel) models.AddressModel {
	address := models.AddressModel{}
	address = originalAddress

	if len(address.StreetAddress) == 0 {
		address.StreetAddress = parsedAddress.StreetAddress
	}
	if len(address.City) == 0 {
		address.City = parsedAddress.City
	}
	if len(address.State) == 0 {
		address.State = parsedAddress.State
	}
	if len(address.PostalCode) == 0 {
		address.PostalCode = parsedAddress.PostalCode
	}

	return address
}

func MergeLot(originalLot models.LotModel, parsedLot models.LotModel) models.LotModel {
	lot := models.LotModel{}
	lot = originalLot

	if lot.LotSize == nil {
		lot.LotSize = parsedLot.LotSize
	}
	lot.SpaceModel = MergeSpace(originalLot.SpaceModel, parsedLot.SpaceModel)

	return lot
}

func MergeStructures(originalStructures []*models.StructureModel, parsedStructures []*models.StructureModel) []*models.StructureModel {
	structures := []*models.StructureModel{}
	structures = originalStructures

	// if there are no structures, just set it to the parsed structure
	if len(structures) == 0 {
		return parsedStructures
	}

	hasMainHouse := false
	// will only need to merge main structure for this parser
	for index, structure := range structures {
		if structure.StructureType == "main_house" {
			hasMainHouse = true
			structure.SpaceModel = MergeSpace(originalStructures[index].SpaceModel, parsedStructures[index].SpaceModel)
			if structure.Size == nil {
				structure.Size = parsedStructures[index].Size
			}
			if structure.Levels == nil {
				structure.Levels = parsedStructures[index].Levels
			}
			structure.Rooms = MergeRooms(originalStructures[index].Rooms, parsedStructures[index].Rooms)
		}
	}

	// if there is no main house add it from the parsed structure and return
	if !hasMainHouse {
		structures = append(structures, parsedStructures[0])
		return structures
	}

	return structures
}

func MergeRooms(originalRooms []*models.RoomModel, parsedRooms []*models.RoomModel) []*models.RoomModel {
	rooms := []*models.RoomModel{}
	rooms = originalRooms

	// if there are no rooms in the homdna, just set it to the parsed rooms
	if len(rooms) == 0 {
		return parsedRooms
	}

	roomsMerged := make([]bool, len(rooms))
	parsedRoomsMerged := make([]bool, len(parsedRooms))

	// check to see if room needs to be merged and not just added
	for i, room := range rooms {
		for j, parsedRoom := range parsedRooms {
			// only merge the room if it hasn't been merged into another room of same type or name
			if room.RoomType == parsedRoom.RoomType && *room.Name == *parsedRoom.Name && !roomsMerged[i] && !parsedRoomsMerged[j] {
				roomsMerged[i] = true
				parsedRoomsMerged[j] = true
				// merge room details here
				room.SpaceModel = MergeSpace(originalRooms[i].SpaceModel, parsedRoom.SpaceModel)
				if room.Level <= 1 {
					room.Level = parsedRoom.Level
				}
			}
		}
	}

	// if room hasn't been merge then add it
	for i, val := range parsedRoomsMerged {
		if !val {
			rooms = append(rooms, parsedRooms[i])
		}
	}

	return rooms
}

func MergeSpace(originalSpace models.SpaceModel, parsedSpace models.SpaceModel) models.SpaceModel {
	space := models.SpaceModel{}
	space = originalSpace

	if space.Id == nil {
		space.Id = parsedSpace.Id
	}
	space.Features = AppendFeatures(space.Features, parsedSpace.Features)
	space.Appliances = AppendAppliances(space.Appliances, parsedSpace.Appliances)

	return space
}

func AppendFeatures(features []*models.FeatureModel, featuresToAdd []*models.FeatureModel) []*models.FeatureModel {
	for _, feature := range featuresToAdd {
		features = append(features, feature)
	}
	return features
}

func AppendAppliances(appliances []*models.ApplianceModel, appliancesToAdd []*models.ApplianceModel) []*models.ApplianceModel {
	for _, appliance := range appliancesToAdd {
		appliances = append(appliances, appliance)
	}
	return appliances
}
