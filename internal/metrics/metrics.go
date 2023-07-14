package metrics

import (
	"github.com/eatplanted/mikrotik-ros-exporter/internal/mikrotik"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	probeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_probe_success",
		Help: "Whether the Mikrotik probe was successful",
	})
)

func CreateRegistryWithMetrics(client mikrotik.Client) (*prometheus.Registry, error) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(probeMetric)

	if err := setHealthMetrics(client, registry); err != nil {
		probeMetric.Set(0)
		return registry, err
	}

	if err := setResourceMetrics(client, registry); err != nil {
		probeMetric.Set(0)
		return registry, err
	}

	return registry, nil
}

func setHealthMetrics(client mikrotik.Client, registry *prometheus.Registry) error {
	healthResult, err := client.GetHealth()
	if err != nil {
		return err
	}

	temperatureMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_health_temperature",
		Help: "System Health Temperature",
	})
	registry.MustRegister(temperatureMetric)
	temperatureMetric.Set(healthResult.Temperature)

	voltageMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_health_voltage",
		Help: "System Health Voltage",
	})
	registry.MustRegister(voltageMetric)
	voltageMetric.Set(healthResult.Voltage)

	return nil
}

func setResourceMetrics(client mikrotik.Client, registry *prometheus.Registry) error {
	resource, err := client.GetResource()
	if err != nil {
		return err
	}

	cpuCountMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_resource_cpu_count",
		Help: "Number of CPUs present on the system. Each core is separate CPU, Intel HT is also separate CPU.",
	})
	registry.MustRegister(cpuCountMetric)
	cpuCountMetric.Set(resource.CpuCount)

	cpuFrequencyMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_resource_cpu_frequency",
		Help: "Current CPU frequency",
	})
	registry.MustRegister(cpuFrequencyMetric)
	cpuFrequencyMetric.Set(resource.CpuFrequency)

	cpuLoadMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_resource_cpu_load",
		Help: "Percentage of used CPU resources. Combines all CPUs.",
	})
	registry.MustRegister(cpuLoadMetric)
	cpuLoadMetric.Set(resource.CpuLoad)

	freeHddSpaceMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_resource_hdd_space_free",
		Help: "Free space on hard drive in bytes",
	})
	registry.MustRegister(freeHddSpaceMetric)
	freeHddSpaceMetric.Set(resource.FreeHddSpace)

	totalHddSpaceMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_resource_hdd_space_total",
		Help: "Size of the hard drive in bytes",
	})
	registry.MustRegister(totalHddSpaceMetric)
	totalHddSpaceMetric.Set(resource.TotalHddSpace)

	freeMemoryMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_resource_memory_free",
		Help: "Unused amount of RAM in bytes",
	})
	registry.MustRegister(freeMemoryMetric)
	freeMemoryMetric.Set(resource.FreeMemory)

	totalMemoryMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_resource_memory_total",
		Help: "Size of the memory in bytes",
	})
	registry.MustRegister(totalMemoryMetric)
	totalMemoryMetric.Set(resource.TotalMemory)

	writeSectorsSinceRebootMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_resource_write_sectors_since_reboot",
		Help: "writeSectSinceRebootMetric",
	})
	registry.MustRegister(writeSectorsSinceRebootMetric)
	writeSectorsSinceRebootMetric.Set(resource.WriteSectSinceReboot)

	writeSectorsTotalMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mikrotik_system_resource_write_sectors_total",
		Help: "writeSectSinceRebootMetric",
	})
	registry.MustRegister(writeSectorsTotalMetric)
	writeSectorsTotalMetric.Set(resource.WriteSectTotal)

	probeMetric.Set(1)
	return nil
}
