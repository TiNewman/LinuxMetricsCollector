import processListStyles from '../styles/Process_List.module.css'

const processListView = props => (
  <div>
   <h1 className={processListStyles.h1}> Process List </h1>
   <div className="overflow-x-auto flex flex-col justify-center items-center">
     <table class="table">
       <thead>
         <tr>
           <th>PID</th>
           <th>Name</th>
           <th>CPU</th>
           <th>RAM</th>
           <th>Disk</th>
           <th>Status</th>
         </tr>
       </thead>
       <tbody>
         <tr className="hover">
           <th>1</th>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
         </tr>
         <tr className="hover">
           <th>2</th>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
         </tr>
         <tr className="hover">
           <th>3</th>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
         </tr>
         <tr className="hover">
           <th>4</th>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
         </tr>
         <tr className="hover">
           <th>5</th>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
           <td></td>
         </tr>
       </tbody>
     </table>
   </div>
  </div>
);

export default processListView;