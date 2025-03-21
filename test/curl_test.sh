#!/bin/bash

# 测试配置
SERVER_URL="http://localhost:8080"
USER_ID="1001"  # 请替换为实际存在的用户ID
DEVICE_ID="DEV123456"  # 请替换为实际存在的设备ID

# 颜色配置
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 创建临时目录
TEMP_DIR=$(mktemp -d)
echo -e "${BLUE}临时目录: ${TEMP_DIR}${NC}"

# 创建测试证书文件
cat > ${TEMP_DIR}/test_cert.pem << 'EOL'
-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUECPZAZ4aIZnKzuhc9fE4/BUNdnIwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMzA0MTAwODM5NTlaFw0yNDA0
MDkwODM5NTlaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQC8YobzCNZs0jF8G1X3GhJm3HkLZ2s9F5X8O5zQFG+e
88JxCJYgLPYdd9yaiaCSXvYRWMtDcQgDmA3S7YUCh5p1U/dUibMzjnCOgZ7QQS1K
vY7qwTFEHgL+Qo6mW8UAcwraTTY7Avpj5j+Mp3pEEGQzHS7H94PnzPBqpUmIcYgj
2zVbHvLkaIxuxaLGT+a6QMLtZbYQKNvtlqRRZeGO+Zn7/jpzL6/GwSgXfIb/XaUE
L6d67xxUibNI3KtQYZU9V7iCgwdQ5t9xttO0ZqrOlPv4QUs/t3wZmD+TiDkEj1kC
hHbbcYOqRRmFQo/8UUz7tseTYAEkSyFonXmwGHgzAgMBAAGjUzBRMB0GA1UdDgQW
BBTFVyKGaKCYz4NfJeK5YuT0PCxRUDAfBgNVHSMEGDAWgBTFVyKGaKCYz4NfJeK5
YuT0PCxRUDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCkwRtL
lHH8D9uY9OKkGqPPhGJAZX9kkVX1w2pp/aIJJi8j0IsbjhrdGFUVk8UbQgj+Gub0
/NMGqSqkWQT3aNZ41S9wbZSsYQHnvlY6A7+WArcPQ0PXeFakUjwGYR8hPXPzRTlG
YnBHwxuLCRvO9XBCbN3GbUJjL2vlHwJyqw0QmvAi3xTUm/8pnVnVWGE+W/wFbAnc
5/Vp0mGq8GRxHIz0UYEZVYLNBPzoQoGYbJLYwLN5+jQpAIpZdi03SGPLfznGJ6jY
FuMIPI0NLvkW+7Ur+fkG55bNH0hx8jVcDx+Vqq/GFqvvxN10uGVB3v2MrI3uDEbv
PB4oFHHUdE4TUrfr
-----END CERTIFICATE-----
EOL

# 创建测试密钥文件
cat > ${TEMP_DIR}/test_key.pem << 'EOL'
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC8YobzCNZs0jF8
G1X3GhJm3HkLZ2s9F5X8O5zQFG+e88JxCJYgLPYdd9yaiaCSXvYRWMtDcQgDmA3S
7YUCh5p1U/dUibMzjnCOgZ7QQS1KvY7qwTFEHgL+Qo6mW8UAcwraTTY7Avpj5j+M
p3pEEGQzHS7H94PnzPBqpUmIcYgj2zVbHvLkaIxuxaLGT+a6QMLtZbYQKNvtlqRR
ZeGO+Zn7/jpzL6/GwSgXfIb/XaUEL6d67xxUibNI3KtQYZU9V7iCgwdQ5t9xttO0
ZqrOlPv4QUs/t3wZmD+TiDkEj1kChHbbcYOqRRmFQo/8UUz7tseTYAEkSyFonXmw
GHgzAgMBAAECggEALEf3AnmEQ8vvtTvTWyf7rDpB7xx5BTCxq6qszdyHCJXtVgdX
WFC2Y3kKWj2ATm0SiWHRSOKM7QvjNr8Nb5YdYhcD4mMTX2hkSKFuMviGh5xWnPKS
AqHa3kGPO89gQvQQkrJJh6GUhhaGXAarR8jzIkEY+TlVpMRnEyzL/2CY6ZHGS3e3
c2GQqRlczQuFJQOxJJO3Tn86ZVtLS1FE33QdAkXx1JLz5fKFRULvW67ZZ5AHhjrF
bEBZxRXjS0q4hcRNi2AHgKKx/bK6xyCTKFoHjnHp2JMqvBz/FCDFKrPGAHMQvfES
N5GRwEVCjk/bmOxYXcKcRJIYJcO6tSGxSO2BtNJ9AQKBgQD/D0h9QpXyXtyBPg+Z
4q+W4CK/DoXSrOJU1jj7XE7aIxXHR5s2/dx6GQjvwvgZXgHJpGVJBvB/UxGhKb2m
C17xx3r+fWuOZDBXdkv0C9BoeuEUK2qJbuU1PDXV2iuTAi5GEwKDCrHiBfUNP+Nn
lVLdyn/fA8IQwkXIhdEurLTCXwKBgQC9DOCHzHVSzCnDLGgH5jEVbCdLAZEjIrbZ
c96QceSYfwQZqFxnDyDXiwUnUYOV3H1+pgWmQJEe2/9Jf9dMiJGsXxPxNpXCaqGj
Uo1vJgAQyCKQpPfWNcKRnmGLNf5Uy44NQDgYV1gJl+WRjOH+SPRtmcYvNrEqAUbL
VFcKvSTqnQKBgQDQsaeR6a0W0T76lxpUJU1dX/O1iRtgQK6t0Ib+42UPxJnr8OFU
WVb8wgKo7KNVEGskS3JOEiIQPBnqRpW2jJJfHEya60q3o16zC7Rc0BxkqQfTt+2e
kdrndfWVj1Xp45CiaFikG0YhMvwM6JtKM9Uy09xGiJzBAK2+Q3R0vwJYbQKBgE9a
G3MclXYK7sD5mQzPztQu63zJ+OFG9t4XQaxC9oj41HHXbVqevbxLMYgxSsWGa8+o
ZY7dQntF6xsToSu3TsW8tG3aLXELOEU6aKkY2t2JMQAqKJoXrwXfWogC3ytfUnfZ
1qPuoPrq8iqNtmKGxanUtESrpXoAkJwZDjpGLnvtAoGAC3HpgfDQn2ma1eBUiMWf
rh5iQJXgxn99lDh5QxuHSZjQqQXZwYDYYHQfXf+8vE7YC3Bj0rJ4acekpsBz36gn
G1FWmwAAEVu4qvBcC6DDMSctpIbJEyXZOSBvBGxWZ/gL3ooSpCfMaIPKKl3Invrx
QA5BICbW1gtrIByvXbpPuPE=
-----END PRIVATE KEY-----
EOL

echo -e "${BLUE}测试文件已创建：${NC}"
echo -e "${BLUE}- 证书文件: ${TEMP_DIR}/test_cert.pem${NC}"
echo -e "${BLUE}- 密钥文件: ${TEMP_DIR}/test_key.pem${NC}"

# 函数：执行curl命令并格式化输出
function test_api() {
    local title=$1
    local cmd=$2
    
    echo -e "\n${BLUE}========== $title ==========${NC}"
    echo -e "${BLUE}执行命令: ${NC}$cmd"
    
    # 执行命令并保存结果
    local result=$(eval $cmd)
    local status=$?
    
    if [ $status -eq 0 ]; then
        echo -e "${GREEN}成功！${NC}"
        echo "响应:"
        echo "$result" | python -m json.tool 2>/dev/null || echo "$result"
    else
        echo -e "${RED}失败！错误码: $status${NC}"
        echo "$result"
    fi
}

# 测试1：绑定用户证书
test_api "测试用户证书绑定" "curl -s -X POST -F 'cert=@${TEMP_DIR}/test_cert.pem' ${SERVER_URL}/bind/users/${USER_ID}/cert"

# 测试2：绑定用户密钥
test_api "测试用户密钥绑定" "curl -s -X POST -F 'key=@${TEMP_DIR}/test_key.pem' ${SERVER_URL}/bind/users/${USER_ID}/key"

# 测试3：绑定设备证书
test_api "测试设备证书绑定" "curl -s -X POST -F 'cert=@${TEMP_DIR}/test_cert.pem' ${SERVER_URL}/bind/devices/${DEVICE_ID}/cert"

# 测试4：绑定设备密钥
test_api "测试设备密钥绑定" "curl -s -X POST -F 'key=@${TEMP_DIR}/test_key.pem' ${SERVER_URL}/bind/devices/${DEVICE_ID}/key"

# 测试5：获取用户证书信息
test_api "获取用户证书信息" "curl -s '${SERVER_URL}/cert/info?type=user&id=${USER_ID}'"

# 测试6：获取设备证书信息
test_api "获取设备证书信息" "curl -s '${SERVER_URL}/cert/info?type=device&id=${DEVICE_ID}'"

# 清理临时文件
echo -e "\n${BLUE}清理临时文件...${NC}"
rm -rf ${TEMP_DIR}
echo -e "${GREEN}临时文件已删除${NC}"

echo -e "\n${GREEN}所有测试完成！${NC}" 