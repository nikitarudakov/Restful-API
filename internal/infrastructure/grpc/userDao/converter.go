package userDao

import "git.foxminded.ua/foxstudent106092/user-management/internal/business/model"

func UnmarshalToUser(u *User) *model.User {
	return &model.User{
		Username: u.Username,
		Password: u.Password,
		Role:     u.Role,
	}
}

func MarshalTogRPCUserObj(u *model.User) *User {
	return &User{
		Username: u.Username,
		Password: u.Password,
		Role:     u.Role,
	}
}

func MarshalTogRPCUserReplace(u *model.User, oldUsername string) *UserReplace {
	return &UserReplace{
		User:                     MarshalTogRPCUserObj(u),
		ToUpdateUserWithUsername: &Username{Val: oldUsername},
	}
}

func MarshalTogRPCInsertResult(ir *model.InsertResult) *InsertResult {
	var insertResult InsertResult

	insertResult.ObjectID = []byte(ir.Id.Hex())
	insertResult.Username = ir.Username

	return &insertResult
}
