import os, json, sequtils, strutils, times

if existsEnv("STATUS_PORT"):
    const customPort = $getEnv("STATUS_PORT")

    echo " @dd parsed port: ", customPort
    var nodeCfg = %* {
        "KeyStoreDir": 3,
        "ShhextConfig": %* {
        "BandwidthStatsEnabled": true
        },
    }

    # Waku V1 config
    nodeCfg["ListenAddr"] = newJString("0.0.0.0:" & customPort)

    # Waku V2 config
    nodeCfg["WakuV2Config"] = %* {
        "Port": parseInt(customPort),
        "UDPPort": parseInt(customPort),
    }

    echo "@dd Node config: ", nodeCfg
else:
    echo "@dd missing env STATUS_PORT"