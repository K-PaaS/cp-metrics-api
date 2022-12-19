

![image](https://user-images.githubusercontent.com/67575226/117615561-d74f8d80-b1a4-11eb-8dc0-70117f96cdad.png)
컨테이너플랫폼을 구성하는 UserPortal에 해당하는 Repository이며 아래 Repository와 함께 플랫폼의 서비스를 구성합니다.
- [PaaS-TA 컨테이너 플랫폼 포탈](https://github.com/PaaS-TA/container-platform-portal-ui)
- [PaaS-TA 컨테이너 플랫폼 Common-API](https://github.com/PaaS-TA/container-platform-portal-common-api)
- [PaaS-TA 컨테이너 플랫폼 API](https://github.com/PaaS-TA/contianer-platform-portal-api)
- [PaaS-TA 컨테이너 플랫폼 Terraman](https://github.com/PaaS-TA/container-platform-terraman)



---

# PaaS-TA 컨테이너 플랫폼 Metrics-API
PaaS-TA에서 제공하는 컨테이너 플랫폼의 관리클러스터에 대한 Metrics 상태정보를 제공하기 위하여 제공하는 API 입니다.

- [PaaS-TA 컨테이너 플랫폼 Metrics-API](#paas-ta-컨테이너-플랫폼-metrics-api)
  - [시작하기](#시작하기)
    - [컨테이너 플랫폼 Metrics-API 빌드 방법](#컨테이너-플랫폼-metrics-api-빌드-방법)
  - [문서](#문서)
  - [개발 환경](#개발-환경)
  - [라이선스](#라이선스)

## 시작하기
PaaS-TA 컨테이너 플랫폼 Metrics-API가 수행하는 애플리케이션 관리 작업은 다음과 같습니다.

- 컨테이너 플랫폼 Cluster CPU/Memroy 현황 수집
- 컨테이너 플랫폼 Cluster Node CPU / Memrory 현황 수집
- 사용량 조회를 위한 RestAPI 제공

### 컨테이너 플랫폼 Metrics-API 빌드 방법
PaaS-TA 컨테이너 플랫폼 Metrics-API 소스 코드를 활용하여 로컬 환경에서 빌드가 필요한 경우 다음 명령어를 입력합니다.
```
$ go build
```


## 문서
- 컨테이너 플랫폼 활용에 대한 정보는 [PaaS-TA 컨테이너 플랫폼](https://github.com/PaaS-TA/paas-ta-container-platform)을 참조하십시오. 


## 개발 환경
PaaS-TA 컨테이너 플랫폼 WEB USER의 개발 환경은 다음과 같습니다.

| Situation                      | Version |
| ------------------------------ | ------- |
| go                             | 1.18    |
| go-resty                       | 2.7     |
| go-sql-driver/mysql            | 1.6.0   |
| gorilla/mux                    | 1.8.0   |
| swaggo/http-swagger            | 1.3.3   |
| swaggo/swag                    | 1.8.1   |
| hashicorp/vault/api            | 1.7.2   |
| gorm.io/gorm                   | 1.23.8  |


## 라이선스
PaaS-TA 컨테이너 플랫폼 WEB USER는 [Apache-2.0 License](http://www.apache.org/licenses/LICENSE-2.0)를 사용합니다.
