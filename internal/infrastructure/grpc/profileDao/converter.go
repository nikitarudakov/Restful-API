package profileDao

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
)

func UnmarshalTogRPCProfileObj(profileObj *ProfileObj) *model.Profile {
	return &model.Profile{
		Nickname:  profileObj.Nickname,
		FirstName: profileObj.FirstName,
		LastName:  profileObj.LastName,
		UpdatedAt: profileObj.UpdatedAt,
		CreatedAt: profileObj.CreatedAt,
		DeletedAt: profileObj.DeletedAt,
	}
}

func MarshalTogRPCProfileObj(p *model.Profile) *ProfileObj {
	return &ProfileObj{
		Nickname:  p.Nickname,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		UpdatedAt: p.UpdatedAt,
		CreatedAt: p.CreatedAt,
		DeletedAt: p.DeletedAt,
	}
}

func MarshalTogRPCProfileUpdate(p *model.Update) *ProfileUpdate {
	return &ProfileUpdate{
		Nickname:  p.Nickname,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		UpdatedAt: p.UpdatedAt,
	}
}

func UnmarshalToUpdate(profileUpdate *ProfileUpdate) *model.Update {
	updatedAt := profileUpdate.GetUpdatedAt()

	return &model.Update{
		Nickname:  profileUpdate.GetNickname(),
		FirstName: profileUpdate.GetFirstName(),
		LastName:  profileUpdate.GetLastName(),
		UpdatedAt: &updatedAt,
	}
}

func MarshalTogRPCInsertResult(ir *model.InsertResult) *InsertResult {
	var insertResult InsertResult

	insertResult.ObjectID = []byte(ir.Id.Hex())
	insertResult.Username = ir.Username

	return &insertResult
}
