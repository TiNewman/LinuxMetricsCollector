import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import styles from '../styles/Process_List.module.css'

let socket

/**
 * historyView
 *
 * Sends Websocket request when this page is loaded for Historical data. When request is received, it populates two radial
 * charts (CPU and RAM) in the HTML return statement and then displays the content on the screen.
 *
 * Return: HTML that contains a label for the History page and the radial charts for CPU and RAM along with the days the
 * average covers
 */
const historyView = props => {
  //variables to set and store incoming history data
  const [cpuHistory, setCPUHistory] = useState([])
  const [ramHistory, setRAMHistory] = useState([])

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
    //make connection to websocket server and immediately send request for historical data
    socket.onopen = () => {
      socket.send(JSON.stringify({"request": "history"}))
      console.log("sent history request")
    };

    //handles the incoming message from the websocket server
    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      var processJSON = JSON.parse(e.data)
      setCPUHistory([processJSON.history]) //might need to adjust based on the message that comes through
      setRAMHistory([processJSON.history]) //might need to adjust based on the message that comes through
    }

    //closes the websocket connection when the browser tab is closed
    return () => {
      console.log("closing socket")
      socket.send(JSON.stringify({"request": "stop"}))
      socket.close()
    };
  }

    /** THIS IS THE MANUAL TEST DATA FOR CLIENT WEBSOCKETS
    const cpuHistory = [{Usage:37.7905493, StartDate:"04/25/22", EndDate:"04/28/22"}]
    const ramHistory = [{Usage:13.7905493, StartDate:"04/25/22", EndDate:"04/28/22"}]
    */

  //Maps over cpuHistory and ramHistory and returns html that is injected into app view to avoid crashing before data is
  //received
  return (
    <div>
      <h1 className={styles.h1}> Historical Data </h1>
      {cpuHistory.map((item, index) =>
        <div className="float-left ml-48 mt-10 pt-10 pl-10">
          <h2 className={styles.h2}> CPU History Average </h2>
          <div className="block p-5 shadow-lg shadow-primary mb-5">
            <div className="radial-progress text-neutral border-4 border-primary bg-primary"
                 style={{"--value":item.AverageCpuUsage.toFixed(2), "--size":"12rem"}}>{item.AverageCpuUsage.toFixed(2)}%</div>
          </div>
          <h2 className={styles.h2}>Time: {item.Start} - {item.End}</h2>
        </div>
      )}
      {ramHistory.map((item, index) =>
        <div className="float-left mt-10 pt-10 pl-10">
          <h2 className={styles.h2}> RAM History Average</h2>
          <div className="block p-5 shadow-lg shadow-primary mb-5">
            <div className="radial-progress text-neutral border-4 border-primary bg-primary"
                 style={{"--value":item.AverageMemUsage.toFixed(2), "--size":"12rem"}}>{item.AverageMemUsage.toFixed(2)}%</div>
          </div>
          <h2 className={styles.h2}>Time: {item.Start} - {item.End}</h2>
        </div>
      )}
    </div>
  )

}
export default historyView;