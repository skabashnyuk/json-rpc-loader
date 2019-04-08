package main

import (
	"github.com/cenkalti/backoff"
	"github.com/eclipse/che-go-jsonrpc"
	"github.com/eclipse/che-go-jsonrpc/event"
	"github.com/eclipse/che-go-jsonrpc/jsonrpcws"
	"log"
	"time"
)

type Loader struct {
	tunnelStatuses                *jsonrpc.Tunnel
	tunnelLogs                    *jsonrpc.Tunnel
	wsUrlMajor, wsUrlMinor, token string
	bus                           *event.Bus
	runtimeID                     RuntimeID
}

func (loader *Loader) Init(wsUrlMajor, wsUrlMinor, token string) {

	loader.wsUrlMajor = wsUrlMajor
	loader.wsUrlMinor = wsUrlMinor
	loader.token = token
	loader.bus = event.NewBus()
	loader.runtimeID = RuntimeID{
		Workspace:   RandStringRunes(10),
		Environment: RandStringRunes(10),
		OwnerId:     RandStringRunes(10)}
	//connect to server
	loader.initConnections()
	//setup tunnels
	loader.PushStatuses()
	loader.PushLogs()
}

func (loader *Loader) Start()  {

	//start installation
	loader.pubStarting()
	loader.pubStartingInstallation("org.eclipse.che.exec")
	loader.broadcastLogs("org.eclipse.che.exec", "Exec Agent binary is downloaded remotely")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 Exec-agent configuration")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27   Server")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27     - Address: :4412")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27     - Base path: ''")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27   Process executor")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27     - Logs dir: /workspace_logs/exec-agent")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ⇩ Registered HTTPRoutes:")
	loader.broadcastLogs("org.eclipse.che.exec", "")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 Process Routes:")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ Start Process ........................... POST   /process")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ Get Process ............................. GET    /process/:pid")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ Kill Process ............................ DELETE /process/:pid")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ Get Process Logs ........................ GET    /process/:pid/logs")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ Get Processes ........................... GET    /process")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 Exec-Agent WebSocket routes:")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ Connect to Exec-Agent(websocket) ........ GET    /connect")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 Exec-Agent liveness route:")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ Check Exec-Agent liveness ............... GET    /liveness")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ⇩ Registered RPCRoutes:")
	loader.broadcastLogs("org.eclipse.che.exec", "")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 Process Routes:")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ process.start")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ process.kill")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ process.subscribe")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ process.unsubscribe")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ process.updateSubscriber")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ process.getLogs")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ process.getProcess")
	loader.broadcastLogs("org.eclipse.che.exec", "2019/04/04 12:04:27 ✓ process.getProcesses")
	loader.pubInstallationCompleted("org.eclipse.che.exec", InstallerStatusDone)
	loader.pubStartingInstallation("org.eclipse.che.terminal")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30 Terminal-agent configuration")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30   Server")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30     - Address: :4411")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30     - Base path: ''")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30   Terminal")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30     - Slave command: ''")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30     - Activity tracking enabled: false")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30 ")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30 ⇩ Registered HTTPRoutes:")
	loader.broadcastLogs("org.eclipse.che.terminal", "")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30 Terminal routes:")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30 ✓ Connect to pty(websocket) ............... GET    /pty")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30 ")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30 Terminal-Agent liveness route:")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30 ✓ Check Terminal-Agent liveness ........... GET    /liveness")
	loader.broadcastLogs("org.eclipse.che.terminal", "2019/04/04 12:04:30 ")
	loader.pubInstallationCompleted("org.eclipse.che.terminal", InstallerStatusRunning)
	loader.pubStartingInstallation("org.eclipse.che.ws-agent")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "Workspace Agent will be downloaded from Workspace Master")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,300[main]             [INFO ] [o.a.c.s.VersionLoggerListener 89]    - Server version:        Apache Tomcat/8.5.35")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,303[main]             [INFO ] [o.a.c.s.VersionLoggerListener 91]    - Server built:          Nov 3 2018 17:39:20 UTC")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,303[main]             [INFO ] [o.a.c.s.VersionLoggerListener 93]    - Server number:         8.5.35.0")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,304[main]             [INFO ] [o.a.c.s.VersionLoggerListener 95]    - OS Name:               Linux")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,304[main]             [INFO ] [o.a.c.s.VersionLoggerListener 97]    - OS Version:            3.10.0-957.5.1.el7.x86_64")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,304[main]             [INFO ] [o.a.c.s.VersionLoggerListener 99]    - Architecture:          amd64")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,304[main]             [INFO ] [o.a.c.s.VersionLoggerListener 101]   - Java Home:             /usr/lib/jvm/java-8-openjdk-amd64/jre")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,304[main]             [INFO ] [o.a.c.s.VersionLoggerListener 103]   - JVM Version:           1.8.0_171-8u171-b11-0ubuntu0.16.04.1-b11")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,305[main]             [INFO ] [o.a.c.s.VersionLoggerListener 105]   - JVM Vendor:            Oracle Corporation")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,305[main]             [INFO ] [o.a.c.s.VersionLoggerListener 107]   - CATALINA_BASE:         /home/user/che/ws-agent")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,305[main]             [INFO ] [o.a.c.s.VersionLoggerListener 109]   - CATALINA_HOME:         /home/user/che/ws-agent")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,306[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Djava.util.logging.config.file=/home/user/che/ws-agent/conf/logging.properties")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,306[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Djava.util.logging.manager=org.apache.juli.ClassLoaderLogManager")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,306[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -XX:MaxRAM=600m")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,306[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -XX:MaxRAMFraction=1")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,306[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -XX:+UseParallelGC")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,307[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -XX:MinHeapFreeRatio=10")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,307[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -XX:MaxHeapFreeRatio=20")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,309[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -XX:GCTimeRatio=4")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,309[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -XX:AdaptiveSizePolicyWeight=90")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,310[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dsun.zip.disableMemoryMapping=true")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,310[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Xms50m")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,310[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dfile.encoding=UTF8")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,310[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Djava.security.egd=file:/dev/./urandom")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,311[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dche.logs.dir=/workspace_logs/ws-agent")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,311[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dche.logs.level=INFO")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,311[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Djuli-logback.configurationFile=file:/home/user/che/ws-agent/conf/tomcat-logger.xml")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,311[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Djdk.tls.ephemeralDHKeySize=2048")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,312[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Djava.protocol.handler.pkgs=org.apache.catalina.webresources")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,312[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dorg.apache.catalina.security.SecurityListener.UMASK=0022")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,312[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dcom.sun.management.jmxremote")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,312[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dcom.sun.management.jmxremote.ssl=false")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,312[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dcom.sun.management.jmxremote.authenticate=false")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,313[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dche.local.conf.dir=/home/user/che/ws-agent/conf/")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,313[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dignore.endorsed.dirs=")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,313[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dcatalina.base=/home/user/che/ws-agent")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,313[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Dcatalina.home=/home/user/che/ws-agent")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,313[main]             [INFO ] [o.a.c.s.VersionLoggerListener 115]   - Command line argument: -Djava.io.tmpdir=/home/user/che/ws-agent/temp")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,448[main]             [INFO ] [o.a.c.http11.Http11NioProtocol 560]  - Initializing ProtocolHandler [\"http-nio-4401\"]")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,462[main]             [INFO ] [o.a.t.util.net.NioSelectorPool 67]   - Using a shared selector for servlet write/read")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,475[main]             [INFO ] [o.a.catalina.startup.Catalina 649]   - Initialization processed in 517 ms")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,530[main]             [INFO ] [c.m.JmxRemoteLifecycleListener 336]  - The JMX Remote Listener has configured the registry on port [32002] and the server on port [32102] for the [Platform] server")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,530[main]             [INFO ] [o.a.c.core.StandardService 416]      - Starting service [Catalina]")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,530[main]             [INFO ] [o.a.c.core.StandardEngine 259]       - Starting Servlet Engine: Apache Tomcat/8.5.35")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:35,746[ost-startStop-1]  [INFO ] [o.a.c.startup.HostConfig 957]        - Deploying web application archive [/home/user/che/ws-agent/webapps/ROOT.war]")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "Downloading java LS")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "writing start script to /home/user/che/ls-java/launch.sh")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:42,064[ost-startStop-1]  [INFO ] [i.WorkspaceProjectSynchronizer 66]   - Workspace ID: workspaceay130l6hm4ajaq4p")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:42,069[ost-startStop-1]  [INFO ] [i.WorkspaceProjectSynchronizer 67]   - API Endpoint: http://che-eclipse-che.192.168.64.67.nip.io/api")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:42,513[rcherInitThread]  [INFO ] [o.e.c.a.s.s.i.LuceneSearcher 159]    - Initial indexing complete after 5 msec ")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:43,029[ost-startStop-1]  [INFO ] [o.a.c.startup.HostConfig 1020]       - Deployment of web application archive [/home/user/che/ws-agent/webapps/ROOT.war] has finished in [7,282] ms")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:43,033[main]             [INFO ] [o.a.c.http11.Http11NioProtocol 588]  - Starting ProtocolHandler [\"http-nio-4401\"]")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:04:43,063[main]             [INFO ] [o.a.catalina.startup.Catalina 700]   - Server startup in 7587 ms")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019/04/04 12:04:53 Start new terminal.")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:00,929[ication Handler]  [INFO ] [j.l.JavaLanguageServerLauncher 123]  - Starting: Init...")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:00,934[ication Handler]  [INFO ] [j.l.JavaLanguageServerLauncher 123]  - Starting: 0% Starting Java Language Server ")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:00,993[rverInitializer]  [WARN ] [o.e.l.j.s.GenericEndpoint 171]       - Unsupported notification method: json/schemaAssociations")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:01,042[rverInitializer]  [WARN ] [o.e.l.j.s.GenericEndpoint 171]       - Unsupported notification method: json/schemaAssociations")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:01,047[rverInitializer]  [INFO ] [.a.l.LanguageServerInitializer 222]  - Initialized language server 'org.eclipse.che.plugin.java.languageserver'")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:01,414[ication Handler]  [INFO ] [j.l.JavaLanguageServerLauncher 123]  - Starting: 20% Starting Java Language Server ")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:04,102[ication Handler]  [INFO ] [j.l.JavaLanguageServerLauncher 123]  - Starting: 35% Starting Java Language Server ")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:04,112[ication Handler]  [INFO ] [j.l.JavaLanguageServerLauncher 123]  - Starting: 100% Starting Java Language Server ")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:04,112[ication Handler]  [INFO ] [j.l.JavaLanguageServerLauncher 123]  - Starting: 100% Starting Java Language Server ")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:04,112[ication Handler]  [INFO ] [j.l.JavaLanguageServerLauncher 123]  - Started: Ready")
	loader.broadcastLogs("org.eclipse.che.ws-agent", "2019-04-04 12:05:04,230[ication Handler]  [INFO ] [j.l.JavaLanguageServerLauncher 123]  - Starting: 100% Starting Java Language Server ")
	loader.pubInstallationCompleted("org.eclipse.che.ws-agent", InstallerStatusRunning)
	loader.pubStartingInstallation("org.eclipse.che.ls.java")
	loader.pubInstallationCompleted("org.eclipse.che.ls.java", InstallerStatusRunning)
	loader.pubBootstrappingDone()

}

func (loader *Loader) initConnections() {
	loader.tunnelStatuses = connectRetryOrFail(loader.wsUrlMajor, loader.token)
	loader.tunnelLogs = connectRetryOrFail(loader.wsUrlMinor, loader.token)
}

func (loader *Loader) Close() {
	loader.tunnelStatuses.Close()
	loader.tunnelLogs.Close()
}

// PushLogs sets given tunnel as consumer of installer logs.
// Connector is used to reconnect to jsonrpc endpoint if
// established connection behind given tunnel was lost.
func (loader *Loader) PushLogs() {
	loader.bus.Sub(&tunnelBroadcaster{
		tunnel: loader.tunnelLogs,
		//reconnectPeriod: 1,
		//reconnectOnce:   &sync.Once{},
	}, InstallerLogEventType)
}

// PushStatuses sets given tunnel as consumer of installer/bootstrapper statuses.
func (loader *Loader) PushStatuses() {
	loader.bus.SubAny(&tunnelBroadcaster{tunnel: loader.tunnelStatuses}, InstallerStatusChangedEventType, StatusChangedEventType)
}

func (loader *Loader) pubStartingInstallation(installer string) {
	loader.bus.Pub(&InstallerStatusChangedEvent{
		Status:    InstallerStatusStarting,
		Installer: installer,
		MachineEvent: MachineEvent{
			MachineName: "dev-machine",
			RuntimeID:   loader.runtimeID,
			Time:        time.Now(),
		},
	})
}

func (loader *Loader) pubInstallationCompleted(installer string, status string) {

	loader.bus.Pub(&InstallerStatusChangedEvent{
		Status:    status,
		Installer: installer,
		MachineEvent: MachineEvent{
			MachineName: "dev-machine",
			RuntimeID:   loader.runtimeID,
			Time:        time.Now(),
		},
	})
}

func (loader *Loader) pubInstallerStatusChangedEvent(event *InstallerStatusChangedEvent) {
	loader.bus.Pub(event)
}

func (loader *Loader) pubStarting() {
	loader.bus.Pub(&StatusChangedEvent{
		Status: StatusStarting,
		MachineEvent: MachineEvent{
			MachineName: "dev-machine",
			RuntimeID:   loader.runtimeID,
			Time:        time.Now(),
		},
	})
}

func (loader *Loader) pubBootstrappingDone() {
	loader.bus.Pub(&StatusChangedEvent{
		Status: StatusDone,
		MachineEvent: MachineEvent{
			MachineName: "dev-machine",
			RuntimeID:   loader.runtimeID,
			Time:        time.Now(),
		},
	})
}

func (loader *Loader) broadcastLogs(installer, text string) {
	loader.bus.Pub(&InstallerLogEvent{
		Stream:    StdoutStream,
		Text:      text,
		Installer: installer,
		MachineEvent: MachineEvent{
			MachineName: "dev-machine",
			RuntimeID:   loader.runtimeID,
			Time:        time.Now(),
		},
	})
}

func connectRetryOrFail(endpoint string, token string) *jsonrpc.Tunnel {
	var result *jsonrpc.Tunnel
	err := backoff.Retry(func() error {
		tun, err2 := connect(endpoint, token)
		if err2 != nil {
			return err2
		} else {
			result = tun
			return nil
		}

	}, backoff.NewExponentialBackOff())
	if err != nil {
		log.Panicf("Couldn't connect to endpoint '%s', due to error '%s'", endpoint, err)
	}

	return result
}

func connect(endpoint string, token string) (*jsonrpc.Tunnel, error) {
	conn, err := jsonrpcws.Dial(endpoint, token)
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewManagedTunnel(conn), nil
}

// Connector encloses implementation specific jsonrpc connection establishment.
type Connector interface {
	Connect() (*jsonrpc.Tunnel, error)
}

type tunnelBroadcaster struct {
	tunnel *jsonrpc.Tunnel
	//reconnectPeriod time.Duration
	//reconnectOnce   *sync.Once
	//loader          *Loader
}

func (tb *tunnelBroadcaster) Accept(e event.E) {
	if err := tb.tunnel.Notify(e.Type(), e); err != nil {
		log.Printf("Trying to send event of type '%s' to closed tunnel '%s'", e.Type(), tb.tunnel.ID())
	}
}

func (tb *tunnelBroadcaster) IsDone() bool {
	return tb.tunnel.IsClosed()
}

func (tb *tunnelBroadcaster) Close() { tb.tunnel.Close() }
