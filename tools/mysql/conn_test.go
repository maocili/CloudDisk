package mysql

import (
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func TestDBConn(t *testing.T) {
	tests := []struct {
		name string
		want *xorm.Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DBConn(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBConn() = %v, want %v", got, tt.want)
			}
		})
	}
}
