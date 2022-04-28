import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import styles from '../styles/Process_List.module.css'

let socket

const historyView = props => {
  /*const [cpuHistory, setCPUHistory] = useState([])
  const [ramHistory, setRAMHistory] = useState([])

    useEffect(() => {
        socket = new WebSocket("ws://localhost:8080/ws");
        socketInitializer()
        return () => {
          console.log("closing socket")
          socket.send(JSON.stringify({"request": "stop"}))
          socket.close()
        };
       }, [])

    const socketInitializer = async () => {
      socket.onopen = () => {
        socket.send(JSON.stringify({"request": "history"}))
        console.log("sent history request")
      };

      socket.onmessage = (e) => {
        console.log("Received Message!: " + e.data)
        var processJSON = JSON.parse(e.data)
        setCPUHistory([processJSON.cpu]) //might need to adjust based on the message that comes through
        setRAMHistory([processJSON.memory]) //might need to adjust based on the message that comes through
      }

      return () => {
        console.log("closing socket")
        socket.send(JSON.stringify({"request": "stop"}))
        socket.close()
      };
    } */

    /** THIS IS THE MANUAL TEST DATA FOR CLIENT WEBSOCKETS */
    const cpuHistory = [{Usage:37.7905493, StartDate:"04/25/22", EndDate:"04/28/22"}]
    const ramHistory = [{Usage:13.7905493, StartDate:"04/25/22", EndDate:"04/28/22"}]
    /**/

  return (
    <div>
      <h1 className={styles.h1}> Historical Data </h1>
      {cpuHistory.map((item, index) =>
        <div className="float-left ml-96 mt-10 pt-10 pl-10">
          <h2 className={styles.h2}> CPU History Average </h2>
          <div className="block p-5 shadow-lg shadow-primary mb-5">
            <div className="radial-progress text-neutral border-4 border-primary bg-primary" style={{"--value":item.Usage.toFixed(2), "--size":"12rem"}}>{item.Usage.toFixed(2)}%</div>
          </div>
          <h2 className={styles.h2}>Time: {item.StartDate} - {item.EndDate}</h2>
        </div>
      )}
      {ramHistory.map((item, index) =>
        <div className="float-left mt-10 pt-10 pl-10">
          <h2 className={styles.h2}> RAM History Average</h2>
          <div className="block p-5 shadow-lg shadow-primary mb-5">
            <div className="radial-progress text-neutral border-4 border-primary bg-primary" style={{"--value":item.Usage.toFixed(2), "--size":"12rem"}}>{item.Usage.toFixed(2)}%</div>
          </div>
          <h2 className={styles.h2}>Time: {item.StartDate} - {item.EndDate}</h2>
        </div>
      )}
    </div>
  )

}
export default historyView;