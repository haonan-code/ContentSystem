package dao

import (
	"contentsystem/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func connDB() *gorm.DB {
	mysqlDB, err := gorm.Open(mysql.Open("root:20000406@tcp(localhost:3306)/?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	db, err := mysqlDB.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	return mysqlDB
}

func TestContentDao_Create(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		detail model.ContentDetail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "内容插入",
			fields: fields{
				connDB(),
			},
			args: args{
				detail: model.ContentDetail{
					Title: "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ContentDao{
				db: tt.fields.db,
			}
			if err := c.Create(tt.args.detail); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
