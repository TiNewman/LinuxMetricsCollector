import Link from 'next/link'

const NavBar = () => (
  <div>
    <div className="navbar bg-base-100">
      <a className="btn btn-ghost normal-case text-xl">Linux Metrics Collector</a>
    </div>
    <ul class="menu bg-secondary w-56 p-2 rounded-box">
      <li><Link href="../pages/index"><a>Dashboard</a></Link></li>
      <li><Link href="../pages/process_list"><a>Process List</a></Link></li>
    </ul>
  </div>
);

export default NavBar;