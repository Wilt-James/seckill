package repositories

import (
	"database/sql"
	"seckill/common"
	"seckill/datamodels"
	"strconv"
)

type IOrder interface {
	Conn() error
	Insert(*datamodels.Order)(int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64)(*datamodels.Order, error)
	SelectAll()([]*datamodels.Order, error)
	SelectAllWithInfo()(map[int]map[string]string, error)
}

type OrderManager struct {
	table string
	mysqlConn *sql.DB
}

func NewOrderManager(table string, db *sql.DB) IOrder {
	return &OrderManager{table:table, mysqlConn:db}
}

// 数据连接
func (o *OrderManager) Conn()(err error) {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "orders"
	}
	return
}

// 插入
func (o *OrderManager) Insert(order *datamodels.Order) (productId int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	sql := "INSERT " + o.table + " SET userID = ?, productID = ?, orderStatus = ?"
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// 删除
func (o *OrderManager) Delete(orderID int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}
	sql := "DELETE FROM " + o.table + " WHERE ID = ?"
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(orderID)
	if err != nil {
		return false
	}
	return true
}

// 更新
func (o *OrderManager) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}

	sql := "UPDATE " + o.table + " SET userID = ?, productID = ?, orderStatus = ? WHERE ID = " + strconv.FormatInt(order.ID, 10)
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return err
	}
	return nil
}

// 根据商品ID查询商品
func (o *OrderManager) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}
	sql := "SELECT * FROM " + o.table + " WHERE ID = " + strconv.FormatInt(orderID, 10)
	row, err := o.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.Order{}, err
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, nil
	}

	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

// 获取所有商品
func (o *OrderManager) SelectAll() (orderArray []*datamodels.Order, err error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT * FROM " + o.table
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}

	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderManager) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT orders.ID, products.productName, orders.orderStatus FROM orders LEFT OUTER JOIN products ON orders.productID = products.ID"
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return common.GetResultRows(rows), nil
}
