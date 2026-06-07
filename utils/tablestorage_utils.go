package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func InitClient() *tablestore.TimeseriesClient {
	// yourInstanceName 填写您的实例名称
	instanceName := "embedded-compt"
	// yourEndpoint 填写您的实例访问地址
	endpoint := "https://embedded-compt.cn-hangzhou.ots.aliyuncs.com"
	// 获取环境变量里的 AccessKey ID 和 AccessKey Secret
	accessKeyId := os.Getenv("TABLESTORE_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("TABLESTORE_ACCESS_KEY_SECRET")

	// 初始化表格存储客户端
	client := tablestore.NewTimeseriesClient(endpoint, instanceName, accessKeyId, accessKeySecret)
	return client
}

// RecordEnvData 记录温湿度数据
func RecordEnvData(uuid string, payload string) {
	client := InitClient()

	// {"timestamp":string, "temp":float, "hum":float}
	var json_result struct {
		TimeStamp string  `json:"timestamp"`
		Temp      float64 `json:"temp"`
		Hum       float64 `json:"hum"`
	}
	err := json.Unmarshal([]byte(payload), &json_result)
	if err != nil {
		slog.Error(fmt.Sprintf("Encounter error when decoding json: %v", err))
		slog.Error(fmt.Sprintf("    Original message: %s", payload))
		return
		// TODO
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", json_result.TimeStamp, loc)
	if err != nil {
		slog.Warn(fmt.Sprintf("Time parse fail: %v", err))
		slog.Warn(fmt.Sprintf("    Original time: %s", json_result.TimeStamp))
		slog.Warn("Using server time instead")
		timestamp = time.Now().In(loc)
	}

	// 构造时序数据行 timeseriesRow。
	// timeseriesKey 标识时间线：度量名称、数据源主机和标签。
	timeseriesKey := tablestore.NewTimeseriesKey()
	timeseriesKey.SetMeasurementName("device_env")
	timeseriesKey.SetDataSource(uuid)

	// timeseriesRow 在 timeseriesKey 的基础上关联时间戳和字段值。
	timeseriesRow := tablestore.NewTimeseriesRow(timeseriesKey)
	timeseriesRow.SetTimeInus(timestamp.UnixMicro() / 1e6 * 1e6)
	timeseriesRow.AddField("temperature",
		tablestore.NewColumnValue(tablestore.ColumnType_DOUBLE, json_result.Temp))
	timeseriesRow.AddField("humidity",
		tablestore.NewColumnValue(tablestore.ColumnType_DOUBLE, json_result.Hum))

	// 构造写入时序数据的请求。
	putTimeseriesDataRequest := tablestore.NewPutTimeseriesDataRequest("env")
	putTimeseriesDataRequest.AddTimeseriesRows(timeseriesRow)

	// 调用时序客户端写入时序数据。
	// putTimeseriesDataResponse, err := client.PutTimeseriesData(putTimeseriesDataRequest)
	_, err = client.PutTimeseriesData(putTimeseriesDataRequest)
	if err != nil {
		slog.Error(fmt.Sprintf("Fail to write into the table: %v", err))
		return
	}

	// // 批量写入可能部分成功：写入失败的行通过索引和错误信息返回，已成功写入的行立即提交。
	// if len(putTimeseriesDataResponse.GetFailedRowResults()) > 0 {
	// 	fmt.Println("[Warning]: 时序数据写入完成，部分行写入失败：")
	// 	for i := 0; i < len(putTimeseriesDataResponse.GetFailedRowResults()); i++ {
	// 		FailedRow := putTimeseriesDataResponse.GetFailedRowResults()[i]
	// 		fmt.Println("[Warning]: 失败行索引：", FailedRow.Index, " 错误信息：", FailedRow.Error)
	// 	}
	// } else {
	// 	fmt.Println("[Info]: 时序数据写入成功！RequestId：", putTimeseriesDataResponse.RequestId)
	// }

	slog.Info("Success to write into the table")
}

// RecordColonyData 记录图片和详细结果的存储路径
func RecordColonyData(uuid string,
	plate_id int,
	timestamp time.Time,
	img_path string,
	txt_path string,
	number int) {
	client := InitClient()

	// 构造时序数据行 timeseriesRow。
	// timeseriesKey 标识时间线：度量名称、数据源主机和标签。
	timeseriesKey := tablestore.NewTimeseriesKey()
	timeseriesKey.SetMeasurementName("device_colony")
	timeseriesKey.SetDataSource(uuid)
	timeseriesKey.AddTag("plate_id", strconv.Itoa(plate_id))

	// timeseriesRow 在 timeseriesKey 的基础上关联时间戳和字段值。
	timeseriesRow := tablestore.NewTimeseriesRow(timeseriesKey)
	timeseriesRow.SetTimeInus(timestamp.UnixMicro() / 1e6 * 1e6)
	timeseriesRow.AddField("image",
		tablestore.NewColumnValue(tablestore.ColumnType_STRING, img_path))
	timeseriesRow.AddField("detail",
		tablestore.NewColumnValue(tablestore.ColumnType_STRING, txt_path))
	timeseriesRow.AddField("number",
		tablestore.NewColumnValue(tablestore.ColumnType_INTEGER, number))

	// 构造写入时序数据的请求。
	putTimeseriesDataRequest := tablestore.NewPutTimeseriesDataRequest("colony")
	putTimeseriesDataRequest.AddTimeseriesRows(timeseriesRow)

	// 调用时序客户端写入时序数据。
	// putTimeseriesDataResponse, err := client.PutTimeseriesData(putTimeseriesDataRequest)
	_, err := client.PutTimeseriesData(putTimeseriesDataRequest)
	if err != nil {
		slog.Error(fmt.Sprintf("Fail to write into the table: %v", err))
		return
	}
	slog.Info("Success to write into the table")
}
