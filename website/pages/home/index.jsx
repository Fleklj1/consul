import CallToAction from '@hashicorp/react-call-to-action'
import UseCases from '@hashicorp/react-use-cases'
import BasicHero from '../../components/basic-hero'
import ConsulEnterpriseComparison from '../../components/consul-enterprise-comparison'
import LearnCallout from '../../components/learn-callout'
import CaseStudyCarousel from '../../components/case-study-carousel'
import ProductFeaturesList from '@hashicorp/react-product-features-list'

export default function HomePage() {
  return (
    <div className="p-home">
      <BasicHero
        brand="consul"
        heading="Service Networking Across Any Cloud or Runtime"
        content="Automate network configurations, discover services, and enable secure connectivity across any cloud or runtime."
        links={[
          {
            text: 'Download',
            url: '/downloads',
            type: 'download',
          },
          {
            text: 'Get Started',
            url: 'https://learn.hashicorp.com/consul',
            type: 'outbound',
          },
        ]}
        backgroundImage
      />
      <ProductFeaturesList
        heading="Why Consul?"
        features={[
          {
            title: 'First Class Kubernetes Experience',
            content:
              'Consul provides a Helm chart for a Kubernetes first experience for Service Discovery and Service Mesh use cases.',
            icon: require('./img/why-consul/kubernetes.svg'),
          },
          {
            title: 'Service Mesh Across Runtime Platforms',
            content:
              'Establish a service mesh between Bare Metal, Virtual Machines, and Kubernetes clusters, across any cloud.',
            icon: require('./img/why-consul/service-mesh-runtime.svg'),
          },
          {
            title: 'Dynamic Load Balancing Configurations',
            content:
              'Use Consul to automate service updates to popular load balancers (F5, NGINX, HAProxy) and eliminate manual configuration.',
            icon: require('./img/why-consul/dynamic-load-balancing.svg'),
          },
          {
            title: 'Secure, Multi-Cloud Service Networking',
            content:
              'Secure any service running in any environment. Consul enables users to automate and secure service to service communication.',
            icon: require('./img/why-consul/cloud.svg'),
          },
          {
            title: 'Service Discovery with Health Checking',
            content:
              'Consul enables detecting the deployment of new services, changes to existing ones, and provides real time agent health to reduce downtime.',
            icon: require('./img/why-consul/health.svg'),
          },
          {
            title: 'Robust Ecosystem',
            content:
              'Consul offers support for and integrations with many popular DevOps and Networking tools.',
            icon: require('./img/why-consul/world.svg'),
          },
        ]}
      />
      <CaseStudyCarousel
        title="Trusted by startups and the world’s largest organizations"
        caseStudies={[
          {
            quote:
              'Kubernetes is the 800-pound gorilla of container orchestration, coming with a price tag. So we looked into alternatives - and fell in love with Nomad.',
            caseStudyURL:
              'https://endler.dev/2019/maybe-you-dont-need-kubernetes/',
            person: {
              firstName: 'Matthias',
              lastName: 'Endler',
              photo:
                'https://www.datocms-assets.com/2885/1582163422-matthias-endler.png',
              title: 'Backend Engineer',
            },
            company: {
              name: 'Trivago',
              logo:
                'https://www.datocms-assets.com/2885/1582162145-trivago.svg',
            },
          },
          {
            quote:
              'We have people who are first-time system administrators deploying applications. There is a guy on our team who worked in IT help desk for 8 years - just today he upgraded an entire cluster himself.',
            caseStudyURL: 'https://www.hashicorp.com/case-studies/roblox/',
            person: {
              firstName: 'Rob',
              lastName: 'Cameron',
              photo:
                'https://www.datocms-assets.com/2885/1582180216-rob-cameron.jpeg',
              title: 'Technical Director of Infrastructure',
            },
            company: {
              name: 'Roblox',
              logo:
                'https://www.datocms-assets.com/2885/1582180369-roblox-color.svg',
            },
          },
          {
            quote:
              'Our customers’ jobs are changing constantly. It’s challenging to dynamically predict demand, what types of jobs, and the resource requirements. We found that Nomad excelled in this area.',
            caseStudyURL:
              'https://www.hashicorp.com/resources/nomad-vault-circleci-security-scheduling',
            person: {
              firstName: 'Rob',
              lastName: 'Zuber',
              photo:
                'https://www.datocms-assets.com/2885/1582180618-rob-zuber.jpeg',
              title: 'CTO',
            },
            company: {
              name: 'CircleCI',
              logo:
                'https://www.datocms-assets.com/2885/1582180745-circleci-logo.svg',
            },
          },
          {
            quote:
              'Adopting Nomad did not require us to change our packaging format — we could continue to package Python in Docker and build binaries for the rest of our applications.',
            caseStudyURL:
              'https://medium.com/@copyconstruct/schedulers-kubernetes-and-nomad-b0f2e14a896',
            person: {
              firstName: 'Cindy',
              lastName: 'Sridharan',
              photo:
                'https://www.datocms-assets.com/2885/1582181517-cindy-sridharan.png',
              title: 'Engineer',
            },
            company: {
              name: 'imgix',
              logo: 'https://www.datocms-assets.com/2885/1582181250-imgix.svg',
            },
          },
        ]}
        logoSection={{
          grayBackground: true,
          featuredLogos: [
            {
              companyName: 'Trivago',
              url:
                'https://www.datocms-assets.com/2885/1582162317-trivago-monochromatic.svg',
            },
            {
              companyName: 'Roblox',
              url:
                'https://www.datocms-assets.com/2885/1582180373-roblox-monochrome.svg',
            },
            {
              companyName: 'CircleCI',
              url:
                'https://www.datocms-assets.com/2885/1582180745-circleci-logo.svg',
            },
          ],
        }}
      />
      <div className="use-cases g-grid-container">
        <h2 className="g-type-display-2">Use Cases</h2>
        <UseCases
          items={[
            {
              title: 'Network Middleware Automation',
              description:
                'Reduce burden of manual, ticket-based networking tasks.',
              image: {
                url: require('./img/use-cases/network-middleware-automation.png?url'),
                format: 'png',
              },
              link: {
                title: 'Learn more',
                url: '/use-cases/network-middleware-automation',
              },
            },
            {
              title: 'Multi-Platform Service Mesh',
              description:
                'Secure, modern application networking across any cloud or runtime.',
              image: {
                url: require('./img/use-cases/multi-platform-service-mesh.png?url'),
                format: 'png',
              },
              link: {
                title: 'Learn more',
                url: '/use-cases/multi-platform-service-mesh',
              },
            },
            {
              title: 'Service Discovery & Health Checks',
              description:
                'Enable services to locate other services running in any environment and provide real-time health status.',
              image: {
                url: require('./img/use-cases/service-discovery-and-health-checks.png?url'),
                format: 'png',
              },
              link: {
                title: 'Learn more',
                url: '/use-cases/service-discovery-and-health-checks',
              },
            },
          ]}
        />
      </div>

      <LearnCallout
        headline="Learn the latest Consul skills"
        brand="consul"
        items={[
          {
            title: 'Getting Started',
            category: 'Step-by-Step Guides',
            time: '48 mins',
            link:
              'https://learn.hashicorp.com/consul?track=getting-started#getting-started',
            image: require('./img/learn/getting-started.svg?url'),
          },
          {
            title: 'Run Consul on Kubernetes',
            category: 'Step-by-Step Guides',
            time: '142 mins',
            link:
              'https://learn.hashicorp.com/consul?track=kubernetes#kubernetes',
            image: require('./img/learn/getting-started.svg?url'),
          },
        ]}
      />

      <ConsulEnterpriseComparison />

      <CallToAction
        heading="Ready to get started?"
        content="Consul open source addresses the technical complexity of managing production services by providing a way to discover, automate, secure and connect applications and networking configurations across distributed infrastructure and clouds."
        brand="consul"
        links={[
          {
            text: 'Explore HashiCorp Learn',
            url: 'https://learn.hashicorp.com/consul',
            type: 'outbound',
          },
          {
            text: 'Explore Documentation',
            url: '/docs',
            type: 'inbound',
          },
        ]}
        variant="compact"
        theme="light"
      />
    </div>
  )
}
