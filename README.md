## Related Repositories

<table>
  <tr>
    <td colspan=2 align=center>플랫폼</td>
    <td colspan=2 align=center><a href="https://github.com/K-PaaS/cp-deployment">컨테이너 플랫폼</a></td>
    <td colspan=2 align=center><a href="https://github.com/K-PaaS/sidecar-deployment">사이드카</a></td>
    <td colspan=2 align=center><a href="https://github.com/K-PaaS/ap-deployment">어플리케이션 플랫폼</a></td>
  </tr>
  <tr>
    <td colspan=2 align=center>포털</td>
    <td colspan=2 align=center><a href="https://github.com/K-PaaS/cp-portal-release">CP 포털</a></td>
    <td colspan=2 align=center>-</td>
    <td colspan=2 align=center><a href="https://github.com/K-PaaS/portal-deployment">AP 포털</a></td>
  </tr>
  <tr align=center>
    <td colspan=2 rowspan=9>Component<br>/ 서비스</td>
    <td colspan=2><a href="https://github.com/K-PaaS/cp-portal-common-api">Common API</a></td>
    <td colspan=2>-</td>
    <td colspan=2><a href="https://github.com/K-PaaS/ap-mongodb-shard-release">MongoDB</a></td>
  </tr>
  <tr align=center>
    <td colspan=2><a href="https://github.com/K-PaaS/cp-metrics-api">🚩Metric API</a></td>
    <td colspan=2>  </td>
    <td colspan=2><a href="https://github.com/K-PaaS/ap-mysql-release">MySQL</a></td>
  </tr>
  <tr align=center>
    <td colspan=2><a href="https://github.com/K-PaaS/cp-portal-api">Portal API</a></td>
    <td colspan=2>  </td>
    <td colspan=2><a href="https://github.com/K-PaaS/ap-pipeline-release">Pipeline</a></td>
  </tr>
  <tr align=center>
    <td colspan=2><a href="https://github.com/K-PaaS/cp-portal-ui">Portal UI</a></td>
    <td colspan=2>  </td>
    <td colspan=2><a href="https://github.com/K-PaaS/ap-rabbitmq-release">RabbintMQ</a></td>
  </tr>
  <tr align=center>
    <td colspan=2><a href="https://github.com/K-PaaS/cp-portal-service-broker">Service Broker</a></td>
    <td colspan=2>  </td>
    <td colspan=2><a href="https://github.com/K-PaaS/ap-on-demand-redis-release">Redis</a></td>
  </tr>
  <tr align=center>
    <td colspan=2><a href="https://github.com/K-PaaS/cp-terraman">Terraman API</a></td>
    <td colspan=2>  </td>
    <td colspan=2><a href="https://github.com/K-PaaS/ap-source-control-release">SoureceControl</a></td>
  </tr>
</table>

<i>🚩 You are here.</i>


<br>

## K-PaaS 컨테이너 플랫폼 Metric API
K-PaaS 컨테이너 플랫폼의 관리 클러스터에 대한 Metrics 상태정보를 제공하는 REST API 입니다.
- [시작하기](#시작하기)
   - [컨테이너 플랫폼 Metrics API 빌드 방법](#컨테이너-플랫폼-metrics-api-빌드-방법)
- [문서](#문서)
- [개발 환경](#개발-환경)
- [라이선스](#라이선스)

<br>

## 시작하기
K-PaaS 컨테이너 플랫폼 Metrics API가 수행하는 애플리케이션 관리 작업은 다음과 같습니다.
- 컨테이너 플랫폼 Cluster 리소스 현황 수집
- 컨테이너 플랫폼 Cluster CPU/Memrory 현황 수집
- 사용량 조회를 위한 RestAPI 제공

#### 컨테이너 플랫폼 Metrics API 빌드 방법
K-PaaS 컨테이너 플랫폼 Metrics API 소스 코드를 활용하여 로컬 환경에서 빌드가 필요한 경우 다음 명령어를 입력합니다.
```
$ go build
```

#### 컨테이너 플랫폼 환경에서 컨테이너 이미지 빌드하고 HarborRepository에 업로드하는 방법
```shell
## image build
$ sudo podman build -t harbor.{HarborRepositoryIP}.nip.io/cp-portal-repository/cp-portal-metric-api .

# image push to harbor
$ sudo podman push harbor.{HarborRepositoryIP}.nip.io/cp-portal-repository/cp-portal-metric-api:latest
```

<br>

## 문서
- 컨테이너 플랫폼 활용에 대한 정보는 [K-PaaS 컨테이너 플랫폼](https://github.com/K-PaaS/container-platform)을 참조하십시오.

<br>

## 개발 환경
K-PaaS 컨테이너 플랫폼 Metrics API의 개발 환경은 다음과 같습니다.

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

<br>

## 라이선스
K-PaaS 컨테이너 플랫폼 Metrics API는 [Apache-2.0 License](http://www.apache.org/licenses/LICENSE-2.0)를 사용합니다.
