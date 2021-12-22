include .env

# 特定の Realm をエクスポートする。
# Keycloak コンテナが起動している状態で実行し、コンソールに `Admin console listening on http://127.0.0.1:10090` が表示されたら Ctrl + C で停止させる。
# エクスポートする Realm は .env、エクスポート先は odcker-compose.yaml を参照のこと。
export_keycloak:
	docker-compose exec keycloak /opt/jboss/keycloak/bin/standalone.sh \
		-Djboss.socket.binding.port-offset=100 \
		-Dkeycloak.migration.action=export \
		-Dkeycloak.migration.provider=singleFile \
		-Dkeycloak.migration.realmName=${APP_NAME} \
		-Dkeycloak.migration.usersExportStrategy=REALM_FILE \
		-Dkeycloak.migration.file=/data/realm_${APP_NAME}.json

.PHONY: ¥
	export-keycloak
