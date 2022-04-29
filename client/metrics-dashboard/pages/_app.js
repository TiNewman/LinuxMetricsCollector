import Layout from '../components/Layout'
import '../styles/globals.css'

/**
 * MyApp
 *
 * Entrypoint for the application. All of the pages and other components get loaded into this view
 *
 * Return: the layout component wrapped around any page the user wants to display at a given time
 */
function MyApp({ Component, pageProps }) {
  return (
    <Layout>
      <Component {...pageProps} />
    </Layout>
  )
}

export default MyApp
