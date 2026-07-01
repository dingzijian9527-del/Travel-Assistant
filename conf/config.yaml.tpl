app:
  name: travel-assistant
  env: dev

http:
  host: 0.0.0.0
  port: 8080
  read_timeout: 5s
  write_timeout: 10s

rpc:
  host: 0.0.0.0
  user:
    service_name: user-service
    port: 9001
    target: 127.0.0.1:9001
  ai_agent:
    service_name: ai-agent-service
    port: 9002
    target: 127.0.0.1:9002

log:
  level: info
  dir: logs
  max_age_days: 7
  console: true

mysql:
  dsn: "${TRAVEL_ASSISTANT_MYSQL_DSN}"
  max_idle_conns: 10
  max_open_conns: 50
  conn_max_lifetime_seconds: 3600

redis:
  addr: "${TRAVEL_ASSISTANT_REDIS_ADDR}"
  username: "${TRAVEL_ASSISTANT_REDIS_USERNAME}"
  password: "${TRAVEL_ASSISTANT_REDIS_PASSWORD}"
  db: 0

auth:
  jwt_secret: "${TRAVEL_ASSISTANT_AUTH_JWT_SECRET}"
  jwt_expire: 24h

upload:
  local_dir: uploads
  max_size_mb: 20
  allowed_extensions:
    - .jpg
    - .jpeg
    - .png
    - .webp
    - .pdf

rag:
  enabled: true
  provider: local
  address: "${TRAVEL_ASSISTANT_RAG_ADDRESS}"
  collection_name: travel_knowledge
  embedding_dim: 768
  top_k: 3
  min_score: 0.15

ai:
  provider: ark
  api_key: "${TRAVEL_ASSISTANT_AI_API_KEY}"
  base_url: https://ark.cn-beijing.volces.com/api/v3
  endpoint_id: "${TRAVEL_ASSISTANT_AI_ENDPOINT_ID}"
  model_name: "${TRAVEL_ASSISTANT_AI_MODEL_NAME}"
  model: ""
  timeout: 60s
  stream: true
  system_prompt: |
    你是旅行助手项目中的专属旅行规划智能体，只服务旅游出行场景。
    你必须像多轮对话助手一样理解上下文：用户前面提到过的目的地、天数、预算、同行人、偏好，后续追问时要继续沿用，不要每次都当成新问题。
    回答前先判断用户当前意图：行程规划、天气穿搭、美食推荐、住宿选址、交通路线、景点攻略、预算拆分、避坑提醒、方案修改、非旅游问题。
    不要因为命中某个关键词就输出固定答案，必须结合上下文、目的地、天数、季节、人数、预算和用户追问意图生成差异化回复。
    只有用户明确要求行程、路线、几日游、规划、安排时，才输出行程方案；天气、美食、住宿、交通、预算等问题必须分别回答对应内容，不能统一套旅行计划模板，也不要使用固定标题清单硬套所有问题。
    如果用户问天气但缺少目的地或日期，要追问目的地和日期；如果上下文已有目的地，就直接围绕该目的地给出穿搭、雨具、防晒、行程调整等具体建议。
    如果用户问美食，要给出本地菜、小吃、用餐区域、避坑方式和如何放进行程；如果问住宿，要给出选址逻辑、适合人群、预算差异和避坑点；如果问交通，要给出大交通、市内交通、跨区通勤和返程缓冲建议。
    回复必须具体、可执行，避免空泛套话。不要编造实时天气、实时价格、营业状态；需要实时信息时提醒用户以官方或平台实时信息为准。
    用户说继续、再详细点、换一个、优化、不要太赶时，要基于上一轮方案修改或展开，不要重复生成相同内容。
    非旅游问题统一礼貌婉拒，不延伸回答。
  max_prompt_chars: 2000
