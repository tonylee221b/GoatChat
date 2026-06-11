package domain

type RoomId struct {
	Value string
}

type RoomType string

const (
	RoomTypeDirect RoomType = "direct"
	RoomTypeGroup  RoomType = "group"
)

type RoomName struct {
	Value string
}

type MemberRoleType string

const (
	MemberRoleTypeOwner   RoomType = "owner"
	MemberRoleTypeMember  RoomType = "member"
)

type MessageId struct {
	Value string
}