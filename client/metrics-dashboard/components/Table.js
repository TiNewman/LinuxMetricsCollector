import TableScrollbar from 'react-table-scrollbar'

const Table = ({ data, column }) => (
  <div>
     <div className="overflow-x-auto flex flex-col p-5 justify-center table-zebra items-center">
       <TableScrollbar height="500px">
       <table className="table">
         <thead>
           <tr>
             {column.map((item, index) => <TableHeadItem key={item.heading} item={item}/>)}
           </tr>
         </thead>
         <tbody>
           {data.map((item, index) => <TableRow key={item.PID} item={item} column={column} />)}
         </tbody>
       </table>
       </TableScrollbar>
     </div>
   </div>
)

const TableHeadItem = ({ item }) => <th>{item.heading}</th>
const TableRow = ({ item, column }) => (
  <tr>
    {column.map((columnItem, index) => {
      return <td key={item.PID + columnItem.value}>{item[`${columnItem.value}`]}</td>
    })}
  </tr>
)

export default Table;