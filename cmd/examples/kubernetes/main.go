package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/kubernetes"
)

func randomString() string {
	return strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(rand.Int64(), 10)
}

func main() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}

	wg := sync.WaitGroup{}

	c := client.NewMgcClient(apiToken, client.WithRetryConfig(15, 2*time.Second, 60*time.Second, 2.0))
	k8sClient := kubernetes.New(c)

	deleteAllClusters(k8sClient)

	var idComNodePool string
	var idSemNodePool string

	fmt.Println("Creating clusters")
	wg.Add(1)
	go func() {
		idComNodePool = ""
		fmt.Println("Creating cluster with node pool")
		idComNodePool = ExampleCreateCluster(k8sClient)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		idSemNodePool = ""
		fmt.Println("Creating cluster without node pool")
		idSemNodePool = ExampleCreateClusterWithoutNodepool(k8sClient)
		wg.Done()
	}()

	wg.Wait()
	time.Sleep(10 * time.Second)
	fmt.Println("idComNodePool", idComNodePool)
	fmt.Println("idSemNodePool", idSemNodePool)

	ExampleListClusters(k8sClient)

	fmt.Println("Waiting for clusters to be ready")
	wg.Add(1)
	go func() {
		WaitClusterRunning(k8sClient, idComNodePool) // ok
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		WaitClusterRunning(k8sClient, idSemNodePool) // ok
		wg.Done()
	}()

	wg.Wait()

	ExampleGetCluster(k8sClient, idComNodePool)    // ok
	ExampleGetKubeConfig(k8sClient, idComNodePool) // ok
	ExampleUpdateCluster(k8sClient, idComNodePool) // ok

	idNodepool := ExampleGetNodePoolsList(k8sClient, idComNodePool)
	ExampleGetNodePool(k8sClient, idComNodePool, idNodepool)
	ExampleNodePoolOperations(k8sClient, idComNodePool)
	ExampleNodePoolOperationsWithEmptyTaints(k8sClient, idComNodePool)
	ExampleNodePoolOperationsWithTaints(k8sClient, idComNodePool)

	ExampleListFlavorsAndVersions(k8sClient)
	ExampleDeleteCluster(k8sClient, idSemNodePool)
	ExampleDeleteCluster(k8sClient, idComNodePool)
}

func deleteAllClusters(k8sClient *kubernetes.KubernetesClient) {
	clusters, err := k8sClient.Clusters().List(context.Background(), kubernetes.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, cluster := range clusters {
		err = k8sClient.Clusters().Delete(context.Background(), cluster.ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Cluster deleted:", cluster.ID)
	}

	//Check if all clusters are deleted
	for {
		clusters, err = k8sClient.Clusters().List(context.Background(), kubernetes.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}

		if len(clusters) == 0 {
			break
		}

		time.Sleep(10 * time.Second)
		fmt.Println("Waiting for clusters to be deleted")
	}

	fmt.Println("All clusters deleted")
}

func WaitClusterRunning(k8sClient *kubernetes.KubernetesClient, clusterID string) {
	for {
		cluster, err := k8sClient.Clusters().Get(context.Background(), clusterID, nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(cluster.Status.State)

		if cluster.Status.State == "Running" {
			break
		}

		time.Sleep(10 * time.Second)
	}

	fmt.Println("Cluster running")
}

func ExampleCreateClusterWithoutNodepool(k8sClient *kubernetes.KubernetesClient) string {

	// Criar um novo cluster
	createReq := kubernetes.ClusterRequest{
		Name:         randomString(),
		Version:      "v1.30.2",
		Description:  "Cluster de exemplo",
		NodePools:    []kubernetes.CreateNodePoolRequest{},
		AllowedCIDRs: []string{"192.168.0.0/24"},
	}

	cluster, err := k8sClient.Clusters().Create(context.Background(), createReq)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Printf("Cluster criado com ID: %s\n", cluster.ID)
	return cluster.ID
}

func ExampleCreateCluster(k8sClient *kubernetes.KubernetesClient) string {

	// Criar um novo cluster
	createReq := kubernetes.ClusterRequest{
		Name:        randomString(),
		Version:     "v1.30.2",
		Description: "Cluster de exemplo",
		NodePools: []kubernetes.CreateNodePoolRequest{
			{
				Name:     randomString(),
				Flavor:   "cloud-k8s.gp1.small",
				Replicas: 3,
			},
		},
		AllowedCIDRs: []string{"192.168.0.0/24"},
	}

	cluster, err := k8sClient.Clusters().Create(context.Background(), createReq)
	if err != nil {
		log.Println(err)
		return ""
	}

	fmt.Printf("Cluster criado com ID: %s\n", cluster.ID)
	return cluster.ID
}

func ExampleListClusters(k8sClient *kubernetes.KubernetesClient) {

	// Listar clusters com paginação
	clusters, err := k8sClient.Clusters().List(context.Background(), kubernetes.ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Expand: []string{"node_pools"},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nClusters listados:")

	for _, cluster := range clusters {
		fmt.Printf("%s - %s (%s) - %s\n", cluster.ID, cluster.Name, cluster.Flavor, cluster.Status.State)
	}
}

func ExampleGetCluster(k8sClient *kubernetes.KubernetesClient, clusterID string) {
	ctx := context.Background()

	// Obter detalhes do cluster
	cluster, err := k8sClient.Clusters().Get(ctx, clusterID, []string{"node_pools"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nDetalhes do Cluster %s:\n", clusterID)
	fmt.Printf(" - Versão: %s\n", cluster.Version)
	fmt.Printf(" - Status: %s\n", cluster.Status.State)
	fmt.Printf(" - Node Pools: %d\n", len(cluster.NodePools))

}

func ExampleUpdateCluster(k8sClient *kubernetes.KubernetesClient, clusterID string) {
	ctx := context.Background()
	// Atualizar CIDRs permitidos
	updateReq := kubernetes.AllowedCIDRsUpdateRequest{
		AllowedCIDRs: []string{"192.168.0.0/24", "10.0.0.0/16"},
	}

	updatedCluster, err := k8sClient.Clusters().Update(ctx, clusterID, updateReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nCIDRs atualizados:", updatedCluster.AllowedCIDRs)
}

func ExampleGetKubeConfig(k8sClient *kubernetes.KubernetesClient, clusterID string) {
	ctx := context.Background()
	kubeconfig, err := k8sClient.Clusters().GetKubeConfig(ctx, clusterID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nKubeconfig:", kubeconfig.CurrentContext)
}

func ExampleNodePoolOperationsWithTaints(k8sClient *kubernetes.KubernetesClient, clusterID string) {

	ctx := context.Background()

	// Criar novo node pool
	poolReq := kubernetes.CreateNodePoolRequest{
		Name:     randomString(),
		Flavor:   "cloud-k8s.gp1.small",
		Replicas: 1,
		Tags:     []string{"ai"},
		Taints: []kubernetes.Taint{
			{
				Key:    "gpu",
				Effect: "NoSchedule",
			},
		},
	}

	newPool, err := k8sClient.Nodepools().Create(ctx, clusterID, poolReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nNode Pool criado: %s (%s)\n", newPool.Name, newPool.ID)

	// Listar node pools
	pools, err := k8sClient.Nodepools().List(ctx, clusterID, kubernetes.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNode Pools:")
	for _, pool := range pools {
		fmt.Printf(" - %s (%d replicas)\n", pool.Name, pool.Replicas)
	}
	pool, err := k8sClient.Nodepools().Get(ctx, clusterID, newPool.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNode Pool:", pool)
	// Deletar node pool
	err = k8sClient.Nodepools().Delete(ctx, clusterID, newPool.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nNode Pool deletado com sucesso")
}
func ExampleNodePoolOperationsWithEmptyTaints(k8sClient *kubernetes.KubernetesClient, clusterID string) {

	ctx := context.Background()

	// Criar novo node pool
	poolReq := kubernetes.CreateNodePoolRequest{
		Name:     randomString(),
		Flavor:   "cloud-k8s.gp1.small",
		Replicas: 1,
		Tags:     []string{"ai"},
		Taints:   []kubernetes.Taint{},
	}

	newPool, err := k8sClient.Nodepools().Create(ctx, clusterID, poolReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nNode Pool criado: %s (%s)\n", newPool.Name, newPool.ID)

	// Listar node pools
	pools, err := k8sClient.Nodepools().List(ctx, clusterID, kubernetes.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNode Pools:")
	for _, pool := range pools {
		fmt.Printf(" - %s (%d replicas)\n", pool.Name, pool.Replicas)
	}

	pool, err := k8sClient.Nodepools().Get(ctx, clusterID, newPool.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNode Pool:", pool)

	// Deletar node pool
	err = k8sClient.Nodepools().Delete(ctx, clusterID, newPool.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nNode Pool deletado com sucesso")
}

func ExampleNodePoolOperations(k8sClient *kubernetes.KubernetesClient, clusterID string) {

	ctx := context.Background()

	// Criar novo node pool
	poolReq := kubernetes.CreateNodePoolRequest{
		Name:     randomString(),
		Flavor:   "cloud-k8s.gp1.small",
		Replicas: 1,
		Tags:     []string{"ai"},
	}

	newPool, err := k8sClient.Nodepools().Create(ctx, clusterID, poolReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nNode Pool criado: %s (%s)\n", newPool.Name, newPool.ID)

	// Listar node pools
	pools, err := k8sClient.Nodepools().List(ctx, clusterID, kubernetes.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNode Pools:")
	for _, pool := range pools {
		fmt.Printf(" - %s (%d replicas)\n", pool.Name, pool.Replicas)
	}

	// Atualizar node pool
	updateReq := kubernetes.PatchNodePoolRequest{
		Replicas: helpers.IntPtr(3),
	}

	updatedPool, err := k8sClient.Nodepools().Update(ctx, clusterID, newPool.ID, updateReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nNode Pool atualizado: %d replicas\n", updatedPool.Replicas)
	pool, err := k8sClient.Nodepools().Get(ctx, clusterID, newPool.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNode Pool:", pool)
	// Deletar node pool
	err = k8sClient.Nodepools().Delete(ctx, clusterID, newPool.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nNode Pool deletado com sucesso")
}

func ExampleListFlavorsAndVersions(k8sClient *kubernetes.KubernetesClient) {

	// Listar flavors disponíveis
	flavors, err := k8sClient.Flavors().List(context.Background(), kubernetes.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nFlavors para Node Pools:")
	for _, f := range flavors.Results[0].ControlPlane {
		fmt.Printf("CP - %s (%d vCPUs, %dMB RAM)\n", f.Name, f.VCPU, f.RAM)
	}
	for _, f := range flavors.Results[0].NodePool {
		fmt.Printf("NP - %s (%d vCPUs, %dMB RAM)\n", f.Name, f.VCPU, f.RAM)
	}
}

func ExampleDeleteCluster(k8sClient *kubernetes.KubernetesClient, clusterID string) {

	// Esperar cluster ficar estável antes de deletar
	err := waitForClusterStatus(context.Background(), k8sClient, clusterID, "active")
	if err != nil {
		log.Fatal(err)
	}

	err = k8sClient.Clusters().Delete(context.Background(), clusterID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nCluster %s deletado com sucesso\n", clusterID)
}

// Helper function para esperar status do cluster
func waitForClusterStatus(ctx context.Context, client *kubernetes.KubernetesClient, clusterID, targetStatus string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout esperando status do cluster")
		case <-ticker.C:
			cluster, err := client.Clusters().Get(ctx, clusterID, nil)
			if err != nil {
				return err
			}

			fmt.Printf("Status atual do cluster: %s\n", cluster.Status.State)

			if cluster.Status.State == targetStatus {
				return nil
			}

			if cluster.Status.State == "error" {
				return fmt.Errorf("cluster em estado de erro: %s", cluster.Status.Message)
			}
		}
	}
}

func ExampleGetNodePoolsList(k8sClient *kubernetes.KubernetesClient, clusterID string) string {

	ctx := context.Background()

	nodePools, err := k8sClient.Nodepools().List(ctx, clusterID, kubernetes.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNode Pools:", nodePools)

	return nodePools[0].ID
}

func ExampleGetNodePool(k8sClient *kubernetes.KubernetesClient, clusterID string, nodePoolID string) {
	ctx := context.Background()

	nodePool, err := k8sClient.Nodepools().Get(ctx, clusterID, nodePoolID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNode Pool:", nodePool)
}
