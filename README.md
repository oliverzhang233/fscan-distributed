
# Fscan 分布式扫描系统  

## 项目概述  
本项目是一个基于 **Fscan** 的分布式扫描系统，通过在子公司内网部署 **Agent 探针** 实现扫描任务的下发与结果收集，最终在总部服务器集中查看扫描结果。系统支持任务分发、分布式扫描、结果汇总等功能，具备高安全性、可靠性和可扩展性，适用于企业级内网安全扫描场景。  


## 系统架构  
```mermaid
graph TD
    A[总部服务器] -->|下发任务| B{消息队列 RabbitMQ}
    B --> C[子公司 Agent 探针]
    C -->|执行 Fscan 扫描| C
    C -->|上报结果| B
    B --> A[总部服务器]
    A --> D[可视化结果界面]
```  


## 技术选型  
| 组件                | 技术/工具                | 说明                          |  
|---------------------|-------------------------|-------------------------------|  
| 总部服务器语言       | Go                      | 高性能、并发友好，适合服务端开发 |  
| Agent 端语言         | Go                      | 轻量高效，支持多平台部署        |  
| 消息队列             | RabbitMQ                | 任务与结果的异步通信            |  
| 配置文件             | YAML                    | 结构化配置管理                  |  
| 部署脚本             | Shell                   | 一键构建与启动                  |  


## 目录结构  
```plaintext
fscan-system/
├── server/                # 总部服务器代码
│   ├── main.go            # 服务端主逻辑
│   ├── server_config.yaml # 服务端配置文件
│   └── deploy_server.sh   # 服务端部署脚本
├── agent/                 # Agent 探针代码
│   ├── main.go            # Agent 主逻辑
│   ├── agent_config.yaml  # Agent 配置文件
│   └── deploy_agent.sh    # Agent 部署脚本
└── README.md              # 项目说明文档
```  


## 部署指南  


### 前置条件  
1. **安装依赖**  
   - **Go 环境**：`go 1.18+`（[下载 Go](https://go.dev/dl/)）  
   - **RabbitMQ**：消息队列服务（[安装 RabbitMQ](https://www.rabbitmq.com/download.html)）  
   - **Fscan 工具**：从 [Fscan 开源仓库](https://github.com/shadow1ng/fscan) 下载对应平台的可执行文件  


### 一、总部服务器部署  
#### 1. 克隆项目  
```bash
git clone https://github.com/your-username/fscan-system.git
cd fscan-system/server
```  

#### 2. 配置文件 `server_config.yaml`  
```yaml
amqp_url: "amqp://guest:guest@localhost:5672/"  # RabbitMQ 连接地址
task_queue: "scan_tasks"                        # 任务队列名称
result_queue: "scan_results"                    # 结果队列名称
server_port: 8080                               # 服务端端口
```  

#### 3. 一键构建与启动  
```bash
chmod +x deploy_server.sh
./deploy_server.sh
```  
- 服务端启动后监听 `:8080` 端口，提供 `/send_task`（任务下发）和 `/receive_result`（结果接收）接口  


### 二、Agent 探针部署  
#### 1. 克隆项目  
```bash
git clone https://github.com/your-username/fscan-system.git
cd fscan-system/agent
```  

#### 2. 配置文件 `agent_config.yaml`  
```yaml
amqp_url: "amqp://guest:guest@localhost:5672/"  # RabbitMQ 连接地址（需与服务端一致）
task_queue: "scan_tasks"                        # 任务队列名称（需与服务端一致）
result_queue: "scan_results"                    # 结果队列名称（需与服务端一致）
```  

#### 3. 放置 Fscan 可执行文件  
- **方案一**：与 Agent 程序同目录  
  ```plaintext
  agent/
  ├── agent          # Agent 可执行文件（构建后生成）
  ├── fscan          # Fscan 可执行文件（需手动放置）
  └── ...            # 其他文件
  ```  
- **方案二**：系统 `PATH` 目录（如 `/usr/local/bin`）  
  ```bash
  sudo cp /path/to/fscan /usr/local/bin/
  ```  

#### 4. 一键构建与启动  
```bash
chmod +x deploy_agent.sh
./deploy_agent.sh
```  
- Agent 启动后监听 RabbitMQ 任务队列，收到任务后自动执行 Fscan 扫描并上报结果  


## 使用方法  
### 1. 下发扫描任务  
通过 HTTP POST 请求调用服务端接口：  
```bash
curl -X POST http://your-server:8080/send_task \
  -H "Content-Type: application/json" \
  -d '{"host": "192.168.1.1", "args": "-p 1-1000 -u admin -P password.txt"}'
```  
- `host`：目标主机（必填）  
- `args`：Fscan 扫描参数（选填，默认 `-h` 基础扫描）  

### 2. 查看扫描结果  
服务端接收到结果后打印日志，可扩展数据库存储和可视化界面（如添加 Grafana/Elasticsearch 等组件）。  


## 安全与优化建议  
1. **消息队列安全**  
   - 修改 RabbitMQ 默认账号密码（`guest:guest` 仅用于测试）  
   - 使用 SSL/TLS 加密通信数据  

2. **Agent 权限控制**  
   - 限制 Agent 探针的网络访问权限，仅允许与总部服务器和 RabbitMQ 通信  
   - 在企业内网中通过 VPN/专线部署，避免公网暴露  

3. **高可用性**  
   - 部署多个 Agent 探针实现负载均衡  
   - 为 RabbitMQ 和总部服务器添加集群支持  


## 贡献与反馈  
欢迎提交 Issues 反馈问题或建议，也可通过 Pull Requests 参与代码贡献。  
- GitHub 仓库：[your-username/fscan-system](https://github.com/your-username/fscan-system)  


