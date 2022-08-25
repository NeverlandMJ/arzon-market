package user_test

import (
	"fmt"
	"testing"

	"github.com/NeverlandMJ/arzon-market/pkg/user"
	"github.com/stretchr/testify/require"
)

func TestService_CreateUser_EmptyInput(t *testing.T) {
	cases := []struct {
		name string
		args user.PreSignUpUser
		want user.User
		err  error
	}{
		{
			name: "should get error",
			args: user.PreSignUpUser{
				Name:        "",
				PhoneNumber: "887882307",
				Password:    "123",
			},
			want: user.User{},
			err:  fmt.Errorf("empty"),
		},
		{
			name: "should get error",
			args: user.PreSignUpUser{
				Name:        "sunbula",
				PhoneNumber: "",
				Password:    "123",
			},
			want: user.User{},
			err:  fmt.Errorf("empty"),
		},
		{
			name: "should get error",
			args: user.PreSignUpUser{
				Name:        "sunbula",
				PhoneNumber: "887882307",
				Password:    "",
			},
			want: user.User{},
			err:  fmt.Errorf("empty"),
		},
		{
			name: "should get error",
			args: user.PreSignUpUser{
				Name:        "",
				PhoneNumber: "",
				Password:    "",
			},
			want: user.User{},
			err:  fmt.Errorf("empty"),
		},
		{
			name: "should pass",
			args: user.PreSignUpUser{
				Name:        "Hasanova Sunbula",
				Password:    "213",
				PhoneNumber: "1234567",
			},
			want: user.User{
				ID: "",
				FullName: "Hasanova Sunbula",
				Password: "213",
				PhoneNumber: "1234567",
			},
			err:  nil,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			t.Helper()
			got, err := user.NewUser(v.args.Name, v.args.Password, v.args.PhoneNumber)
			v.want.ID = got.ID
			require.Equal(t, v.err, err)
			require.Equal(t, v.want, *got)
		})
	}

}
