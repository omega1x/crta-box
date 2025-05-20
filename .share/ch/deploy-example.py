#!/bin/python3

# DESCRIPTION
#   Example of *ClickHouse*-database deployment
#
#   NOTE: If [DO_REDEPLOY] mode is selected then manually 
#         stop all containers and delete directory [PLATFORM_DATA_PATH]


import logging as put
import os
from socket import gethostname
from subprocess import run

# GLOBAL CONFIGURATION
DO_REDEPLOY = False  # if True it removes all previous deployments and starts the fresh one
HOST_IP = r'127.0.0.1' 
PLATFORM_DATA_PATH = os.path.abspath(r'./platform-data')
SCRIPT_EXIT_CODE = {r'OK' : 0, r'Unsuitable infrastructure': 1, r'System command fail': 2}

# LOCAL FUNCTION LIBRARY


# SCRIPT SETUP

## Setup logging
put.basicConfig(
    handlers=[
        put.FileHandler(f'deploy_{gethostname()}.log'),
        put.StreamHandler()
    ],
    format=r'%(asctime)s | %(levelname)s - %(message)s', level=put.INFO
)
put.info(f'Start execution of [{gethostname()}] server deployment')
put.info(f'All further deployment will be made for host [{HOST_IP}]')
put.info(f'All deployment data will be written to storage [{PLATFORM_DATA_PATH}]')
put.info(f'Deployment script will {"do" if DO_REDEPLOY else "not"} reset Docker-system')


## Check environment
put.info(r'Check suitability of infrastructure...')
is_docker_installed = (run(r'command -v docker', shell=True, check=False).returncode == 0)
if not is_docker_installed:
    put.warning(f'Docker does not installed in the system. Infrastructure cannot be deployed!')
    exit(SCRIPT_EXIT_CODE[r'Unsuitable infrastructure'])
    pass

is_nmap_installed = (run(r'command -v nmap', shell=True, check=False).returncode == 0)
if not is_nmap_installed:
    put.warning(f'[nmap] does not installed in the system. Infrastructure cannot be deployed!')
    exit(SCRIPT_EXIT_CODE[r'Unsuitable infrastructure'])
    pass

is_storage_ready = not(DO_REDEPLOY and os.path.isdir(PLATFORM_DATA_PATH))
if DO_REDEPLOY and (not is_storage_ready):
    put.warning(f'Storage [{PLATFORM_DATA_PATH}] is not ready for (re)deployment. Delete directory [{PLATFORM_DATA_PATH}]')
    exit(SCRIPT_EXIT_CODE[r'Unsuitable infrastructure'])
    pass
put.info(r'Ok. Infrastructure is suitable')

# IaC LOGIC
put.info(r'Start IaC logic...')

## Common setup
### Prune Docker-state
if DO_REDEPLOY:
    put.info(r'Reset Docker system for redeployment...')
    is_docker_pruned = (run(r'docker system prune --all --force --volumes', shell=True, check=False).returncode == 0)
    if not is_docker_pruned:
        put.error(r'Error resetting Docker system. Reinstall it.')
        exit(SCRIPT_EXIT_CODE['System command fail'])
        pass
    put.info(r'Ok. Docker system has been successfully reset for redeployment')

### (Re)create storage
os.makedirs(PLATFORM_DATA_PATH, exist_ok=True)


### Run Click-House [https://clickhouse.com/docs]
CLICKHOUSE_IMAGE_NAME = r'clickhouse/clickhouse-server:25.4'  # set the actual version of the server here
CLICKHOUSE_CONTAINER_NAME = r'clickhouse-server'

put.info(f'Deploy service [{CLICKHOUSE_CONTAINER_NAME}] using [{CLICKHOUSE_IMAGE_NAME}] docker-image ...')

is_service_deployed = (
    (run(f'docker image inspect {CLICKHOUSE_IMAGE_NAME}', shell=True, check=False).returncode == 0) or 
    (run(f'docker container inspect {CLICKHOUSE_CONTAINER_NAME}', shell=True, check=False).returncode == 0)
)    

if is_service_deployed:
    put.warning(f'Ok. Service {CLICKHOUSE_CONTAINER_NAME} is already deployed. Go to next service...')
    pass
else:
    CLICKHOUSE_HOST_PATH  = os.path.join(PLATFORM_DATA_PATH, CLICKHOUSE_CONTAINER_NAME)
    put.info(f'Deploy {CLICKHOUSE_CONTAINER_NAME}-image data to [{CLICKHOUSE_HOST_PATH}]')

    ch_data_path = os.path.join(CLICKHOUSE_HOST_PATH, r'data')
    ch_logs_path = os.path.join(CLICKHOUSE_HOST_PATH, r'log')
    ch_conf_path = os.path.join(CLICKHOUSE_HOST_PATH, r'conf.d')
    ch_user_path = os.path.join(CLICKHOUSE_HOST_PATH, r'users.d')

    put.info(f'Create necessary subdirectories in [{CLICKHOUSE_HOST_PATH}]')
    [os.makedirs(name=dn, exist_ok=True) for dn in [
        CLICKHOUSE_HOST_PATH, ch_data_path, ch_logs_path, ch_conf_path, ch_user_path
      ]
    ]
    cmd = f'\
    docker run -d \
      -e CLICKHOUSE_DB=BOXes \
      -e CLICKHOUSE_USER=user \
      -e CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1 \
      -e CLICKHOUSE_PASSWORD=pass \
       \
      -v {ch_data_path}:/var/lib/clickhouse/ \
      -v {ch_logs_path}:/var/log/clickhouse-server/ \
      -v {ch_conf_path}:/etc/clickhouse-server/config.d/ \
      -v {ch_user_path}:/etc/clickhouse-server/users.d/ \
      \
      -p {HOST_IP}:28123:8123 \
      -p {HOST_IP}:29000:9000 \
      \
      --ulimit nofile=262144:262144 \
      --name {CLICKHOUSE_CONTAINER_NAME} \
      {CLICKHOUSE_IMAGE_NAME} \
    '
    is_service_deployed = (run(cmd, shell=True, check=False).returncode == 0)
    if is_service_deployed:
        put.info(f'Service [{CLICKHOUSE_CONTAINER_NAME}] is successfully deployed and run')
        
        #### Add runtime configuration
        put.info(f'Runtime configuration - add listening capabilities')
        configuration_file = open(os.path.join(ch_conf_path, r'listen.xml'), "w")
        configuration_file.write(
            """
            <clickhouse>
                <!-- Listen wildcard address to allow accepting connections from other containers and host network. -->
                <listen_host>::</listen_host>
                <listen_host>0.0.0.0</listen_host>
                <listen_try>1</listen_try>
            </clickhouse>
           """
        )
        configuration_file.close()
        del(configuration_file)

        #### Add users
        put.info(f'Runtime configuration - add users')
        users_file = open(os.path.join(ch_user_path, r'integration-users.xml'), "w")
        users_file.write(
            """
            <clickhouse>
              <!-- Docs: <https://clickhouse.com/docs/en/operations/settings/settings_users/> -->
              <users>
                <writer>
                  <networks>
                    <ip>::/0</ip>
                  </networks>
                  <password>writer!</password>
                  <grants>
                    <query>GRANT SELECT ON *.*</query>
                    <query>GRANT INSERT ON *.*</query>
                  </grants>
                </writer>
                <reader>
                  <networks>
                    <ip>::/0</ip>
                  </networks>
                  <password>reader!</password>
                  <grants>
                    <query>GRANT SELECT ON *.*</query>
                  </grants>
                </reader>
              </users>
            </clickhouse>
           """
        )
        users_file.close()
        del(users_file)
        
        put.info(f'Ok! Finish deploy service [{CLICKHOUSE_CONTAINER_NAME}]')
        pass
    else:
        put.error(f'Errors are detected while deploying service [{CLICKHOUSE_CONTAINER_NAME}]')
        pass
    pass
put.info(r'All done, thanks!')
