import Header from "./Header";
import NavBar from "./NavBar";
import Head from "next/head";

/**
 * Layout
 *
 * Main component for the web UI. All other pages and components get injected
 *
 * Param: children: react component that is injected to add to main view of the application
 */
const Layout = ({children}) => (
  <div className="Layout">
    <Head>
      <title>Metrics Collector</title>
    </Head>
    <NavBar />
    {children}
  </div>
);

export default Layout;