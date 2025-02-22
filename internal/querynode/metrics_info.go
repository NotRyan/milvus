// Copyright (C) 2019-2020 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License.

package querynode

import (
	"context"
	"os"

	"github.com/milvus-io/milvus/internal/proto/commonpb"
	"github.com/milvus-io/milvus/internal/proto/milvuspb"
	"github.com/milvus-io/milvus/internal/util/metricsinfo"
	"github.com/milvus-io/milvus/internal/util/typeutil"
)

func getSystemInfoMetrics(ctx context.Context, req *milvuspb.GetMetricsRequest, node *QueryNode) (*milvuspb.GetMetricsResponse, error) {
	nodeInfos := metricsinfo.QueryNodeInfos{
		BaseComponentInfos: metricsinfo.BaseComponentInfos{
			Name: metricsinfo.ConstructComponentName(typeutil.QueryNodeRole, Params.QueryNodeID),
			HardwareInfos: metricsinfo.HardwareMetrics{
				IP:           node.session.Address,
				CPUCoreCount: metricsinfo.GetCPUCoreCount(false),
				CPUCoreUsage: metricsinfo.GetCPUUsage(),
				Memory:       metricsinfo.GetMemoryCount(),
				MemoryUsage:  metricsinfo.GetUsedMemoryCount(),
				Disk:         metricsinfo.GetDiskCount(),
				DiskUsage:    metricsinfo.GetDiskUsage(),
			},
			SystemInfo: metricsinfo.DeployMetrics{
				SystemVersion: os.Getenv(metricsinfo.GitCommitEnvKey),
				DeployMode:    os.Getenv(metricsinfo.DeployModeEnvKey),
			},
			// TODO(dragondriver): CreatedTime & UpdatedTime, easy but time-costing
			Type: typeutil.QueryNodeRole,
		},
		SystemConfigurations: metricsinfo.QueryNodeConfiguration{
			SearchReceiveBufSize:         Params.SearchReceiveBufSize,
			SearchPulsarBufSize:          Params.SearchPulsarBufSize,
			SearchResultReceiveBufSize:   Params.SearchResultReceiveBufSize,
			RetrieveReceiveBufSize:       Params.RetrieveReceiveBufSize,
			RetrievePulsarBufSize:        Params.RetrievePulsarBufSize,
			RetrieveResultReceiveBufSize: Params.RetrieveResultReceiveBufSize,

			SimdType: Params.SimdType,
		},
	}
	resp, err := metricsinfo.MarshalComponentInfos(nodeInfos)
	if err != nil {
		return &milvuspb.GetMetricsResponse{
			Status: &commonpb.Status{
				ErrorCode: commonpb.ErrorCode_UnexpectedError,
				Reason:    err.Error(),
			},
			Response:      "",
			ComponentName: metricsinfo.ConstructComponentName(typeutil.QueryNodeRole, Params.QueryNodeID),
		}, nil
	}

	return &milvuspb.GetMetricsResponse{
		Status: &commonpb.Status{
			ErrorCode: commonpb.ErrorCode_Success,
			Reason:    "",
		},
		Response:      resp,
		ComponentName: metricsinfo.ConstructComponentName(typeutil.QueryNodeRole, Params.QueryNodeID),
	}, nil
}
