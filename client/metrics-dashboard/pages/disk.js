import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import processListStyles from '../styles/Process_List.module.css'

let socket

/**
 * diskView
 *
 * Sends Websocket request when this page is loaded for disk data. When request is received, it populates four radial
 * charts in the HTML return statement and then displays the content on the screen.
 *
 * Return: HTML that contains a label for the disk page and the disk radial charts for each disk and information for
 * each disk
 */
const diskView = props => {
  //variables to set and store the incoming disk data
  const [diskUsage, setDiskUsage] = useState([])

  //initialize Websocket and close the connection when this page is unmounted from the view
  useEffect(() => {
      socket = new WebSocket("ws://localhost:8080/ws");
      socketInitializer()
      return () => {
        console.log("closing socket")
        socket.send(JSON.stringify({"request": "stop"}))
        socket.close()
      };
     }, [])

  //initializes methods for the websocket
  const socketInitializer = async () => {
    //make connection to websocket server and immediately send request for disk data
    socket.onopen = () => {
      socket.send(JSON.stringify({"request": "disk"}))
      console.log("sent disk request")
    };

    //handles the incoming message from the websocket server
    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      var processJSON = JSON.parse(e.data)
      //set the incoming data to use in charts
      setDiskUsage(processJSON.disk)
    }

    //closes the websocket connection when the browser tab is closed
    return () => {
      console.log("closing socket")
      socket.send(JSON.stringify({"request": "stop"}))
      socket.close()
    };
  }

  /** THIS IS THE MANUAL TEST DATA FOR CLIENT WEBSOCKETS
  const diskUsage = [{"Name":"/dev/nvme0n1p3","MountPoint":"/","Usage":2.060362882143396,"Size":510405.902336},
                     {"Name":"/dev/nvme0n1p3","MountPoint":"/home","Usage":2.060362882143396,"Size":510405.902336},
                     {"Name":"/dev/nvme0n1p2","MountPoint":"/boot","Usage":25.681494412006668,"Size":1020.70272},
                     {"Name":"/dev/nvme0n1p1","MountPoint":"/boot/efi","Usage":2.3118672372403726,"Size":627.900416}]
  */

  //Maps over the diskUsage and returns html that is injected into app view to avoid crashing before data is received
  return (
    <div>
      <h1 className={processListStyles.h1}> Disk Usage </h1>
        {diskUsage.map((item, index) =>
          <div className="overflow-x-center flex flex-col w-1/6 mx-auto ml-24 mt-5 float-left px-4">
            <h2 className={processListStyles.h2}>{item.Name}</h2>
            <div className="block overflow-x-center flex flex-cols mt-5 p-5 shadow-lg shadow-primary">
              <div className="radial-progress text-neutral border-4 border-primary self-center bg-primary "
                   style={{"--value":item.Usage.toFixed(2), "--size":"10rem"}}>{item.Usage.toFixed(2)}%</div>
            </div>
            <p className="pl-8 pt-8">Mount Point: {item.MountPoint}</p>
            <p className="pl-8">Size: {item.Size.toFixed(2)} MB</p>
          </div>
        )}
    </div>
  )
}

export default diskView;