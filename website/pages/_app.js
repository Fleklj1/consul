import './style.css'
import App from 'next/app'
import NProgress from 'nprogress'
import Router from 'next/router'
import ProductSubnav from '../components/subnav'
import MegaNav from '@hashicorp/react-mega-nav'
import Footer from '../components/footer'
import { ConsentManager, open } from '@hashicorp/react-consent-manager'
import consentManagerConfig from '../lib/consent-manager-config'
import bugsnagClient from '../lib/bugsnag'
import Error from './_error'
import Head from 'next/head'
import HashiHead from '@hashicorp/react-head'

Router.events.on('routeChangeStart', NProgress.start)
Router.events.on('routeChangeError', NProgress.done)
Router.events.on('routeChangeComplete', (url) => {
  setTimeout(() => window.analytics.page(url), 0)
  NProgress.done()
})

// Bugsnag
const ErrorBoundary = bugsnagClient.getPlugin('react')

class NextApp extends App {
  static async getInitialProps({ Component, ctx }) {
    let pageProps = {}

    if (Component.getInitialProps) {
      pageProps = await Component.getInitialProps(ctx)
    } else if (Component.isMDXComponent) {
      // fix for https://github.com/mdx-js/mdx/issues/382
      const mdxLayoutComponent = Component({}).props.originalType
      if (mdxLayoutComponent.getInitialProps) {
        pageProps = await mdxLayoutComponent.getInitialProps(ctx)
      }
    }

    return { pageProps }
  }

  render() {
    const { Component, pageProps } = this.props

    return (
      <ErrorBoundary FallbackComponent={Error}>
        <HashiHead
          is={Head}
          title="Consul by HashiCorp"
          siteName="Consul by HashiCorp"
          description="Consul is a free and open source tool for creating golden images for multiple
          platforms from a single source configuration."
          image="https://www.consul.io/img/og-image.png"
          stylesheet={[
            { href: '/css/nprogress.css' },
            {
              href:
                'https://fonts.googleapis.com/css?family=Open+Sans:300,400,600,700&display=swap',
            },
          ]}
          icon={[{ href: '/favicon.ico' }]}
          preload={[
            { href: '/fonts/klavika/medium.woff2', as: 'font' },
            { href: '/fonts/gilmer/light.woff2', as: 'font' },
            { href: '/fonts/gilmer/regular.woff2', as: 'font' },
            { href: '/fonts/gilmer/medium.woff2', as: 'font' },
            { href: '/fonts/gilmer/bold.woff2', as: 'font' },
            { href: '/fonts/metro-sans/book.woff2', as: 'font' },
            { href: '/fonts/metro-sans/regular.woff2', as: 'font' },
            { href: '/fonts/metro-sans/semi-bold.woff2', as: 'font' },
            { href: '/fonts/metro-sans/bold.woff2', as: 'font' },
            { href: '/fonts/dejavu/mono.woff2', as: 'font' },
          ]}
        />
        <MegaNav product="Consul" />
        <ProductSubnav />
        <div className="content">
          <Component {...pageProps} />
        </div>
        <Footer openConsentManager={open} />
        <ConsentManager {...consentManagerConfig} />
      </ErrorBoundary>
    )
  }
}

export default NextApp
