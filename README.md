## HyperSDS-Provisioner
HyperSDS-Provisioner는 HyperCloud 5.0의 HyperSDS-Operator가 cephadm으로 ceph cluster를 설치, 변경할 때 사용하는 image를 개발하는 프로젝트입니다.

### How To Develop
design에 작성된 예제들을 참고하여 개발하시기 바랍니다.

### Test Prerequisites
1. Root 디렉토리('/')에 파일 r/w 권한이 있어야 합니다.
2. 대상 remote node들에 접근하는 계정으로 SSH 접근이 허용되어야 합니다.
- 만약 remote node로의 userId가 root라면 SSH 설정(/etc/ssh/sshd_config)에서 PermitRootLogin을 Yes로 설정해야 함

3. Root 디렉토리에 `cluster.yaml` 파일이 올바른 포맷으로 존재해야 합니다.
