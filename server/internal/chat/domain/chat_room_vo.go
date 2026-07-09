package domain

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

type RoomId struct {
	Value string
}

func NewRoomId(id string) (RoomId, error) {
	if id == "" {
		return RoomId{}, errors.New("room id is empty")
	}
	_, err := uuid.Parse(id)
	if err != nil {
		slog.Error("UUID parse error room id, ", "error", err.Error())
		return RoomId{}, err
	}

	return RoomId{id}, nil
}

type RoomType string

const (
	RoomTypeUnknown RoomType = "unknown"
	RoomTypeDirect  RoomType = "direct"
	RoomTypeGroup   RoomType = "group"
)

func NewRoomType(rt string) (RoomType, error) {
	if rt != string(RoomTypeDirect) && rt != string(RoomTypeGroup) {
		slog.Info("invalid room type, input room type is : ", "rt", rt)
		return RoomTypeUnknown, errors.New("invalid room type")
	}

	return RoomType(rt), nil
}

type RoomName struct {
	Value string
}

func NewRoomName(rn string) (RoomName, error) {
	if rn == "" {
		slog.Info("invalid room name, input room name is empty")
		return RoomName{}, errors.New("room name is empty")
	}

	return RoomName{rn}, nil
}

type RoomDescription struct {
	Value string
}

const DescriptionMaxSize = 255

func NewRoomDescription(rd string) (RoomDescription, error) {
	if len(rd) > DescriptionMaxSize {
		slog.Info("room description to large")
		return RoomDescription{}, errors.New("room description too many contents")
	}
	return RoomDescription{rd}, nil
}

type RoomOwnerId struct {
	Value string
}

func NewRoomOwnerId(oid string) (RoomOwnerId, error) {
	if oid == "" {
		slog.Info("invalid owner id, input owner id is empty")
		return RoomOwnerId{}, errors.New("owner id is empty")
	}

	if _, err := uuid.Parse(oid); err != nil {
		slog.Info("invalid owner id, input owner id is: ", "ownerId", oid)
		return RoomOwnerId{}, errors.New("invalid owner id")
	}

	return RoomOwnerId{oid}, nil
}

type MemberRoleType string

const (
	MemberRoleTypeUnknown MemberRoleType = "unknwon"
	MemberRoleTypeOwner   MemberRoleType = "owner"
	MemberRoleTypeMember  MemberRoleType = "member"
)

func NewMemberRoleType(rt string) (MemberRoleType, error) {
	if rt != string(MemberRoleTypeOwner) && rt != string(MemberRoleTypeMember) {
		slog.Info("invalid member role type, input role type is : ", "rt", rt)
		return MemberRoleTypeUnknown, errors.New("invalid role type")
	}

	return MemberRoleType(rt), nil
}

type MessageId struct {
	Value string
}

func NewMessageId(id string) (MessageId, error) {
	if id == "" {
		slog.Info("invalid message id, input message id is empty")
		return MessageId{}, errors.New("message id is empty")
	}

	if _, err := uuid.Parse(id); err != nil {
		slog.Info("invalid message id, input message id is: ", "messageId", id)
		return MessageId{}, errors.New("invalid message id")
	}

	return MessageId{id}, nil
}
