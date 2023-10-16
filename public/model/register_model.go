package public_model

import "github.com/Bit-Bridge-Source/BitBridge-AuthService-Go/proto/pb"

type RegisterModel struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (registerModel *RegisterModel) ToCreatedUserRequest() *pb.RegisterRequest {
	return &pb.RegisterRequest{
		Email:    registerModel.Email,
		Username: registerModel.Username,
		Password: registerModel.Password,
	}
}
