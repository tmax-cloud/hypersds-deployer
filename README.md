## HyperSDS-Provisioner
HyperSDS-Provisioner는 HyperCloud 5.0의 HyperSDS-Operator가 cephadm으로 ceph cluster를 설치, 변경할 때 사용하는 Pod의 image를 개발하는 프로젝트입니다.

### How To Build
1. Go binary를 build
	```shell
	make build
	```

2. Build한 binary로 container image를 build
	```shell
	# docker command must be permitted to the account
	make container
	# if you want to use private registry
	make container REGISTRY=1.2.3.4:5000
	```

### Prerequisites
1. 대상 remote node들의 계정으로 SSH 접근이 허용되어야 합니다.
- 만약 remote node로의 userId가 root라면 SSH 설정(/etc/ssh/sshd_config)에서 PermitRootLogin을 Yes로 설정해야 함

2. 대상 remote node의 계정은 Root 디렉토리('/')에 파일 r/w 권한이 있어야 합니다.
- Remote node들의 계정을 root로 사용하는 것을 적극 권장함

3. 대상 remote node에 `/working/config/` 디렉토리가 존재해야 합니다.

4. Container의 `/manifest/` 디렉토리에 `cluster.yaml` 파일이 올바른 형식으로 존재해야 합니다.

5. 해당 Pod의 K8s 환경에 `ceph-conf` 이름의 ConfigMap, `ceph-secret` 이름의 Secret이 존재해야 하며, `data` 필드가 설정돼있지 않아야 합니다.

### How To Develop
design에 작성된 예제들을 참고하여 개발하시기 바랍니다.
