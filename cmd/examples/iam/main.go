package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/iam"
	"gopkg.in/yaml.v3"
)

type AuthFile struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessToken     string `yaml:"access_token"`
	RefreshToken    string `yaml:"refresh_token"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

func main() {
	// Get current CLI profile
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home := os.Getenv("HOME")
		if home != "" {
			dir = path.Join(home, ".config")
		}
	}
	profile, err := os.ReadFile(path.Join(dir, "mgc", "current"))
	if err != nil {
		log.Fatal(err)
	}

	// Get JWT from profile
	authPath := path.Join(dir, "mgc", string(profile), "auth.yaml")
	authFileContent, err := os.ReadFile(authPath)
	if err != nil {
		log.Fatal(err)
	}

	var authFile AuthFile
	err = yaml.Unmarshal(authFileContent, &authFile)
	if err != nil {
		log.Fatal(err)
	}

	c := client.NewMgcClient(client.WithJWToken(authFile.AccessToken))
	iamClient := iam.New(c)

	// Exemplos de uso dos serviços IAM
	ExampleListMembers(iamClient)
	ExampleListRoles(iamClient)
	ExampleListPermissions(iamClient)
	ExampleGetAccessControl(iamClient)
	ExampleListServiceAccounts(iamClient)
	ExampleListScopes(iamClient)

	// Exemplos comentados - descomente para executar operações de escrita
	// memberUUID := ExampleCreateMember(iamClient)
	// ExampleGetMemberGrants(iamClient, memberUUID)
	// ExampleAddMemberGrants(iamClient, memberUUID)
	// ExampleBatchUpdateMembers(iamClient)
	// ExampleDeleteMember(iamClient, memberUUID)

	// roleName := ExampleCreateRole(iamClient)
	// ExampleGetRolePermissions(iamClient, roleName)
	// ExampleEditRolePermissions(iamClient, roleName)
	// ExampleGetRoleMembers(iamClient, roleName)
	// ExampleDeleteRole(iamClient, roleName)

	// saUUID := ExampleCreateServiceAccount(iamClient)
	// ExampleEditServiceAccount(iamClient, saUUID)
	// ExampleServiceAccountAPIKeys(iamClient, saUUID)
	// ExampleDeleteServiceAccount(iamClient, saUUID)

	// ExampleCreateAccessControl(iamClient)
	// ExampleUpdateAccessControl(iamClient)
}

// ExampleListMembers demonstra como listar membros
func ExampleListMembers(iamClient *iam.IAMClient) {
	ctx := context.Background()

	fmt.Println("\n=== Listando Membros ===")
	members, err := iamClient.Members().List(ctx, nil)
	if err != nil {
		log.Printf("Erro ao listar membros: %v", err)
		return
	}

	fmt.Printf("Total de membros: %d\n", len(members))
	for _, member := range members {
		fmt.Printf("  - %s (%s) - %s\n", member.Name, member.Email, member.UUID)
		if member.Profile != nil {
			fmt.Printf("    Profile: %s\n", *member.Profile)
		}
	}

	// Exemplo com filtro por email
	fmt.Println("\n=== Listando Membros com Filtro ===")
	email := "user@example.com"
	members, err = iamClient.Members().List(ctx, &email)
	if err != nil {
		log.Printf("Erro ao listar membros com filtro: %v", err)
		return
	}

	fmt.Printf("Membros encontrados com email %s: %d\n", email, len(members))
}

// ExampleCreateMember demonstra como criar um membro
func ExampleCreateMember(iamClient *iam.IAMClient) string {
	ctx := context.Background()

	fmt.Println("\n=== Criando Membro ===")
	createReq := iam.CreateMember{
		Email:       "newuser@example.com",
		Roles:       []string{"viewer"},
		Permissions: []string{"read:instances"},
	}

	member, err := iamClient.Members().Create(ctx, createReq)
	if err != nil {
		log.Printf("Erro ao criar membro: %v", err)
		return ""
	}

	fmt.Printf("Membro criado com sucesso:\n")
	fmt.Printf("  UUID: %s\n", member.UUID)
	fmt.Printf("  Email: %s\n", member.Email)
	fmt.Printf("  Nome: %s\n", member.Name)

	return member.UUID
}

// ExampleGetMemberGrants demonstra como obter grants de um membro
func ExampleGetMemberGrants(iamClient *iam.IAMClient, memberUUID string) {
	ctx := context.Background()

	fmt.Println("\n=== Obtendo Grants do Membro ===")
	grants, err := iamClient.Members().Grants().Get(ctx, memberUUID)
	if err != nil {
		log.Printf("Erro ao obter grants: %v", err)
		return
	}

	fmt.Printf("Roles: %v\n", grants.Roles)
	fmt.Printf("Permissions: %v\n", grants.Permissions)
}

// ExampleAddMemberGrants demonstra como adicionar/remover grants de um membro
func ExampleAddMemberGrants(iamClient *iam.IAMClient, memberUUID string) {
	ctx := context.Background()

	fmt.Println("\n=== Adicionando Grants ao Membro ===")
	addGrant := iam.EditGrant{
		Operation: iam.OperationAdd,
		Roles:     []string{"admin"},
	}

	err := iamClient.Members().Grants().Add(ctx, memberUUID, addGrant)
	if err != nil {
		log.Printf("Erro ao adicionar grants: %v", err)
		return
	}

	fmt.Println("Grants adicionados com sucesso")

	// Exemplo de remoção
	fmt.Println("\n=== Removendo Grants do Membro ===")
	removeGrant := iam.EditGrant{
		Operation: iam.OperationRemove,
		Roles:     []string{"viewer"},
	}

	err = iamClient.Members().Grants().Add(ctx, memberUUID, removeGrant)
	if err != nil {
		log.Printf("Erro ao remover grants: %v", err)
		return
	}

	fmt.Println("Grants removidos com sucesso")
}

// ExampleBatchUpdateMembers demonstra atualização em lote de membros
func ExampleBatchUpdateMembers(iamClient *iam.IAMClient) {
	ctx := context.Background()

	fmt.Println("\n=== Atualização em Lote de Membros ===")
	batchReq := iam.BatchUpdateMembers{
		MemberIDs:       []string{"uuid1", "uuid2"},
		Operation:       iam.OperationAdd,
		RoleNames:       []string{"admin"},
		PermissionNames: []string{"read:instances"},
	}

	err := iamClient.Members().Grants().BatchUpdate(ctx, batchReq)
	if err != nil {
		log.Printf("Erro na atualização em lote: %v", err)
		return
	}

	fmt.Println("Membros atualizados em lote com sucesso")
}

// ExampleDeleteMember demonstra como deletar um membro
func ExampleDeleteMember(iamClient *iam.IAMClient, memberUUID string) {
	ctx := context.Background()

	fmt.Println("\n=== Deletando Membro ===")
	err := iamClient.Members().Delete(ctx, memberUUID)
	if err != nil {
		log.Printf("Erro ao deletar membro: %v", err)
		return
	}

	fmt.Printf("Membro %s deletado com sucesso\n", memberUUID)
}

// ExampleListRoles demonstra como listar roles
func ExampleListRoles(iamClient *iam.IAMClient) {
	ctx := context.Background()

	fmt.Println("\n=== Listando Roles ===")
	roles, err := iamClient.Roles().List(ctx, nil)
	if err != nil {
		log.Printf("Erro ao listar roles: %v", err)
		return
	}

	fmt.Printf("Total de roles: %d\n", len(roles))
	for _, role := range roles {
		fmt.Printf("  - %s (origin: %s)\n", role.Name, role.Origin)
		if role.Description != nil {
			fmt.Printf("    Description: %s\n", *role.Description)
		}
	}

	// Exemplo com filtro por nome
	fmt.Println("\n=== Listando Roles com Filtro ===")
	roleName := "admin"
	roles, err = iamClient.Roles().List(ctx, &roleName)
	if err != nil {
		log.Printf("Erro ao listar roles com filtro: %v", err)
		return
	}

	fmt.Printf("Roles encontradas com nome %s: %d\n", roleName, len(roles))
}

// ExampleCreateRole demonstra como criar uma role
func ExampleCreateRole(iamClient *iam.IAMClient) string {
	ctx := context.Background()

	fmt.Println("\n=== Criando Role ===")
	createReq := iam.CreateRole{
		Name:        "custom-role",
		Description: strPtr("Role customizada para exemplo"),
		Permissions: []string{"read:instances", "read:networks"},
	}

	roles, err := iamClient.Roles().Create(ctx, createReq)
	if err != nil {
		log.Printf("Erro ao criar role: %v", err)
		return ""
	}

	if len(roles) > 0 {
		fmt.Printf("Role criada com sucesso: %s\n", roles[0].Name)
		return roles[0].Name
	}

	return ""
}

// ExampleGetRolePermissions demonstra como obter permissões de uma role
func ExampleGetRolePermissions(iamClient *iam.IAMClient, roleName string) {
	ctx := context.Background()

	fmt.Println("\n=== Obtendo Permissões da Role ===")
	rolePerms, err := iamClient.Roles().Permissions(ctx, roleName)
	if err != nil {
		log.Printf("Erro ao obter permissões: %v", err)
		return
	}

	fmt.Printf("Role: %s\n", rolePerms.Name)
	fmt.Printf("Origin: %s\n", rolePerms.Origin)
	if rolePerms.Description != nil {
		fmt.Printf("Description: %s\n", *rolePerms.Description)
	}
	fmt.Printf("Permissions: %v\n", rolePerms.Permissions)
}

// ExampleEditRolePermissions demonstra como editar permissões de uma role
func ExampleEditRolePermissions(iamClient *iam.IAMClient, roleName string) {
	ctx := context.Background()

	fmt.Println("\n=== Editando Permissões da Role ===")
	editReq := iam.EditPermissions{
		Add:    []string{"write:instances"},
		Remove: []string{"read:networks"},
	}

	roles, err := iamClient.Roles().EditPermissions(ctx, roleName, editReq)
	if err != nil {
		log.Printf("Erro ao editar permissões: %v", err)
		return
	}

	fmt.Printf("Permissões editadas com sucesso. Roles atualizadas: %d\n", len(roles))
}

// ExampleGetRoleMembers demonstra como obter membros de uma role
func ExampleGetRoleMembers(iamClient *iam.IAMClient, roleName string) {
	ctx := context.Background()

	fmt.Println("\n=== Obtendo Membros da Role ===")
	members, err := iamClient.Roles().Members(ctx, roleName)
	if err != nil {
		log.Printf("Erro ao obter membros: %v", err)
		return
	}

	fmt.Printf("Membros com role %s: %d\n", roleName, len(members))
	for _, member := range members {
		fmt.Printf("  - UUID: %s, Roles: %v\n", member.MemberUUID, member.Roles)
	}
}

// ExampleDeleteRole demonstra como deletar uma role
func ExampleDeleteRole(iamClient *iam.IAMClient, roleName string) {
	ctx := context.Background()

	fmt.Println("\n=== Deletando Role ===")
	err := iamClient.Roles().Delete(ctx, roleName)
	if err != nil {
		log.Printf("Erro ao deletar role: %v", err)
		return
	}

	fmt.Printf("Role %s deletada com sucesso\n", roleName)
}

// ExampleListPermissions demonstra como listar produtos e permissões
func ExampleListPermissions(iamClient *iam.IAMClient) {
	ctx := context.Background()

	fmt.Println("\n=== Listando Produtos e Permissões ===")
	products, err := iamClient.Permissions().ProductsAndPermissions(ctx, nil)
	if err != nil {
		log.Printf("Erro ao listar produtos e permissões: %v", err)
		return
	}

	fmt.Printf("Total de produtos: %d\n", len(products))
	for _, product := range products {
		fmt.Printf("  - %s (%d permissões)\n", product.Name, len(product.Permissions))
		for _, perm := range product.Permissions {
			fmt.Printf("    * %s: %s\n", perm.Name, perm.Description)
		}
	}

	// Exemplo com filtro por produto
	fmt.Println("\n=== Listando Permissões com Filtro ===")
	productName := "compute"
	products, err = iamClient.Permissions().ProductsAndPermissions(ctx, &productName)
	if err != nil {
		log.Printf("Erro ao listar permissões com filtro: %v", err)
		return
	}

	fmt.Printf("Permissões do produto %s: %d\n", productName, len(products))
}

// ExampleGetAccessControl demonstra como obter configuração de controle de acesso
func ExampleGetAccessControl(iamClient *iam.IAMClient) {
	ctx := context.Background()

	fmt.Println("\n=== Obtendo Configuração de Controle de Acesso ===")
	ac, err := iamClient.AccessControl().Get(ctx)
	if err != nil {
		log.Printf("Erro ao obter controle de acesso: %v", err)
		return
	}

	if ac.Name != nil {
		fmt.Printf("Nome: %s\n", *ac.Name)
	}
	if ac.Description != nil {
		fmt.Printf("Description: %s\n", *ac.Description)
	}
	fmt.Printf("Enabled: %v\n", ac.Enabled)
	fmt.Printf("Enforce MFA: %v\n", ac.EnforceMFA)
}

// ExampleCreateAccessControl demonstra como criar configuração de controle de acesso
func ExampleCreateAccessControl(iamClient *iam.IAMClient) {
	ctx := context.Background()

	fmt.Println("\n=== Criando Configuração de Controle de Acesso ===")
	createReq := iam.AccessControlCreate{
		Name:        strPtr("custom-ac"),
		Description: strPtr("Configuração customizada"),
	}

	ac, err := iamClient.AccessControl().Create(ctx, createReq)
	if err != nil {
		log.Printf("Erro ao criar controle de acesso: %v", err)
		return
	}

	fmt.Printf("Controle de acesso criado com sucesso\n")
	if ac.Name != nil {
		fmt.Printf("Nome: %s\n", *ac.Name)
	}
}

// ExampleUpdateAccessControl demonstra como atualizar configuração de controle de acesso
func ExampleUpdateAccessControl(iamClient *iam.IAMClient) {
	ctx := context.Background()

	fmt.Println("\n=== Atualizando Configuração de Controle de Acesso ===")
	updateReq := iam.AccessControlStatus{
		Status:     boolPtr(true),
		EnforceMFA: boolPtr(true),
	}

	ac, err := iamClient.AccessControl().Update(ctx, updateReq)
	if err != nil {
		log.Printf("Erro ao atualizar controle de acesso: %v", err)
		return
	}

	fmt.Printf("Controle de acesso atualizado com sucesso\n")
	fmt.Printf("Enabled: %v\n", ac.Enabled)
	fmt.Printf("Enforce MFA: %v\n", ac.EnforceMFA)
}

// ExampleListServiceAccounts demonstra como listar service accounts
func ExampleListServiceAccounts(iamClient *iam.IAMClient) {
	ctx := context.Background()

	fmt.Println("\n=== Listando Service Accounts ===")
	sas, err := iamClient.ServiceAccounts().List(ctx)
	if err != nil {
		log.Printf("Erro ao listar service accounts: %v", err)
		return
	}

	fmt.Printf("Total de service accounts: %d\n", len(sas))
	for _, sa := range sas {
		fmt.Printf("  - %s (%s) - %s\n", sa.Name, sa.Email, sa.UUID)
		if sa.Description != nil {
			fmt.Printf("    Description: %s\n", *sa.Description)
		}
		fmt.Printf("    Tenant: %s\n", sa.Tenant.LegalName)
	}
}

// ExampleCreateServiceAccount demonstra como criar um service account
func ExampleCreateServiceAccount(iamClient *iam.IAMClient) string {
	ctx := context.Background()

	fmt.Println("\n=== Criando Service Account ===")
	createReq := iam.ServiceAccountCreate{
		Name:        "my-service-account",
		Description: "Service account para exemplo",
		Email:       "sa-example@example.com",
	}

	sa, err := iamClient.ServiceAccounts().Create(ctx, createReq)
	if err != nil {
		log.Printf("Erro ao criar service account: %v", err)
		return ""
	}

	fmt.Printf("Service account criado com sucesso:\n")
	fmt.Printf("  UUID: %s\n", sa.UUID)
	fmt.Printf("  Nome: %s\n", sa.Name)
	fmt.Printf("  Email: %s\n", sa.Email)

	return sa.UUID
}

// ExampleEditServiceAccount demonstra como editar um service account
func ExampleEditServiceAccount(iamClient *iam.IAMClient, saUUID string) {
	ctx := context.Background()

	fmt.Println("\n=== Editando Service Account ===")
	editReq := iam.ServiceAccountEdit{
		Name:        strPtr("updated-service-account"),
		Description: strPtr("Descrição atualizada"),
	}

	sa, err := iamClient.ServiceAccounts().Edit(ctx, saUUID, editReq)
	if err != nil {
		log.Printf("Erro ao editar service account: %v", err)
		return
	}

	fmt.Printf("Service account editado com sucesso:\n")
	fmt.Printf("  Nome: %s\n", sa.Name)
}

// ExampleServiceAccountAPIKeys demonstra operações com API keys de service accounts
func ExampleServiceAccountAPIKeys(iamClient *iam.IAMClient, saUUID string) {
	ctx := context.Background()

	fmt.Println("\n=== Listando API Keys do Service Account ===")
	apiKeys, err := iamClient.ServiceAccounts().APIKeys(ctx, saUUID)
	if err != nil {
		log.Printf("Erro ao listar API keys: %v", err)
		return
	}

	fmt.Printf("Total de API keys: %d\n", len(apiKeys))
	for _, key := range apiKeys {
		fmt.Printf("  - %s (%s)\n", *key.Name, key.UUID)
		if key.Description != nil {
			fmt.Printf("    Description: %s\n", *key.Description)
		}
		fmt.Printf("    Scopes: %v\n", key.Scopes)
	}

	// Criar nova API key
	fmt.Println("\n=== Criando API Key ===")
	createKeyReq := iam.APIKeyServiceAccountCreate{
		Name:        "my-api-key",
		Description: strPtr("API key para exemplo"),
		Scopes:      []string{"read:instances", "read:networks"},
	}

	newKey, err := iamClient.ServiceAccounts().CreateAPIKey(ctx, saUUID, createKeyReq)
	if err != nil {
		log.Printf("Erro ao criar API key: %v", err)
		return
	}

	fmt.Printf("API key criada com sucesso:\n")
	fmt.Printf("  UUID: %s\n", newKey.UUID)
	if newKey.KeyPairID != nil {
		fmt.Printf("  Key Pair ID: %s\n", *newKey.KeyPairID)
	}
	if newKey.APIKey != nil {
		fmt.Println("  API Key: [REDACTED]")
	}

	// Editar API key
	fmt.Println("\n=== Editando API Key ===")
	editKeyReq := iam.APIKeyServiceAccountEditInput{
		Name:        strPtr("updated-api-key"),
		Description: strPtr("Descrição atualizada"),
		Scopes:      []string{"read:instances", "write:instances"},
	}

	updatedKey, err := iamClient.ServiceAccounts().EditAPIKey(ctx, saUUID, newKey.UUID, editKeyReq)
	if err != nil {
		log.Printf("Erro ao editar API key: %v", err)
		return
	}

	fmt.Printf("API key editada com sucesso: %s\n", updatedKey.UUID)

	// Revogar API key
	fmt.Println("\n=== Revogando API Key ===")
	err = iamClient.ServiceAccounts().RevokeAPIKey(ctx, saUUID, newKey.UUID)
	if err != nil {
		log.Printf("Erro ao revogar API key: %v", err)
		return
	}

	fmt.Println("API key revogada com sucesso")
}

// ExampleDeleteServiceAccount demonstra como deletar um service account
func ExampleDeleteServiceAccount(iamClient *iam.IAMClient, saUUID string) {
	ctx := context.Background()

	fmt.Println("\n=== Deletando Service Account ===")
	err := iamClient.ServiceAccounts().Delete(ctx, saUUID)
	if err != nil {
		log.Printf("Erro ao deletar service account: %v", err)
		return
	}

	fmt.Printf("Service account %s deletado com sucesso\n", saUUID)
}

// ExampleListScopes demonstra como listar grupos, produtos e scopes
func ExampleListScopes(iamClient *iam.IAMClient) {
	ctx := context.Background()

	fmt.Println("\n=== Listando Grupos, Produtos e Scopes ===")
	groups, err := iamClient.Scopes().GroupsAndProductsAndScopes(ctx)
	if err != nil {
		log.Printf("Erro ao listar scopes: %v", err)
		return
	}

	fmt.Printf("Total de grupos: %d\n", len(groups))
	for _, group := range groups {
		fmt.Printf("  - %s (%s)\n", group.Name, group.UUID)
		fmt.Printf("    Produtos: %d\n", len(group.APIProducts))
		for _, product := range group.APIProducts {
			fmt.Printf("      * %s (%s) - %d scopes\n", product.Name, product.UUID, len(product.Scopes))
			for _, scope := range product.Scopes {
				fmt.Printf("        - %s: %s\n", scope.Name, scope.Title)
			}
		}
	}
}

// Funções auxiliares
func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
