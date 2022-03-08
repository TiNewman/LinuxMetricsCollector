import Link from 'next/link'

const NavBar = () => (
  <div>
    <div className="navbar bg-base-100">
      <a className="btn btn-ghost normal-case text-xl">Linux Metrics Collector</a>
    </div>
    <ul className="menu bg-secondary w-56 p-2 rounded-box">
      <li><Link href="/">Dashboard</Link></li>
      <li><Link href="/process_list">Process List</Link></li>
    </ul>
  </div>
);

export default NavBar;