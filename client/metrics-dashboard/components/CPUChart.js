import Chart from "react-apexcharts"
import { useCallback, useEffect, useState } from 'react'

const CPUChart = (data) => {
  //might need to be in page with websocket code
  const [state,setState] = useState({
                                     options: {
                                       chart: {
                                         id: "basic-bar"
                                       },
                                       xaxis: {
                                         categories: [1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998]
                                       }
                                     },
                                     series: [
                                       {
                                         name: "series-1",
                                         data: [30, 40, 45, 50, 49, 60, 70, 91]
                                       }
                                     ]
                               });
  setState({
                                                options: {
                                                  chart: {
                                                    id: "basic-bar"
                                                  },
                                                  xaxis: {
                                                    categories: [1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998]
                                                  }
                                                },
                                                series: [
                                                  {
                                                    name: "series-1",
                                                    data: data
                                                  }
                                                ]
                                          })

  return (
  <div className="mixed-chart">
    {(typeof window !== 'undefined') &&
    <div className="app">
      <div className="row">
        <div className="mixed-chart">
          <Chart
            options={state.options}
            series={state.series}
            type="bar"
            width="500"
          />
        </div>
      </div>
    </div>
    }
  </div>
  )
}

export default CPUChart