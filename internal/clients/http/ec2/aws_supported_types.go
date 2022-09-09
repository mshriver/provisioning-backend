package ec2

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func NewInstanceTypes(ctx context.Context, types []types.InstanceTypeInfo) ([]*clients.InstanceType, error) {
	list := make([]*clients.InstanceType, 0, len(types))
	for i := range types {
		architectures := types[i].ProcessorInfo.SupportedArchitectures
		for _, arch := range architectures {
			arch, err := clients.MapArchitectures(ctx, string(arch))
			if err != nil {
				return nil, payloads.ClientError(ctx, "Instance type", "", err, 500)
			}

			it := clients.InstanceType{
				Name:         clients.InstanceTypeName(types[i].InstanceType),
				VCPUs:        *types[i].VCpuInfo.DefaultVCpus,
				Cores:        *types[i].VCpuInfo.DefaultCores,
				MemoryMiB:    *types[i].MemoryInfo.SizeInMiB,
				Architecture: arch,
				Supported:    clients.IsSupported(string(types[i].InstanceType)),
			}
			if types[i].InstanceStorageInfo != nil {
				it.EphemeralStorageGB = *types[i].InstanceStorageInfo.TotalSizeInGB
			}
			list = append(list, &it)
		}
	}
	return list, nil
}