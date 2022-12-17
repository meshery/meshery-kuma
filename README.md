<p style="text-align:center;" align="center"><a href="https://layer5.io/meshery"><picture align="center">
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/meshery/meshery-kuma/master/.github/assets/images/meshery/meshery-logo-dark-text-side.svg"  width="70%" align="center" style="margin-bottom:20px;">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/meshery/meshery-kuma/master/.github/assets/images/meshery/meshery-logo-light-text-side.svg" width="70%" align="center" style="margin-bottom:20px;">
  <img alt="Shows an illustrated light mode meshery logo in light color mode and a dark mode meshery logo dark color mode." src="https://raw.githubusercontent.com/meshery/meshery/master/.github/assets/images/meshery/meshery-logo-tag-light-text-side.png" width="70%" align="center" style="margin-bottom:20px;"><br /><br /></p>
  
<p align="center">
<a href="https://hub.docker.com/r/layer5/meshery-kuma" alt="Docker pulls">
  <img src="https://img.shields.io/docker/pulls/layer5/meshery-kuma.svg" /></a>
<a href="https://goreportcard.com/report/github.com/layer5io/meshery-kuma" alt="Go Report Card">
  <img src="https://goreportcard.com/badge/github.com/layer5io/meshery-kuma" /></a>
<a href="(https://github.com/layer5io/meshery-kuma/actions" alt="Build Status">
  <img src="https://github.com/layer5io/meshery-kuma/workflows/Meshery-Kuma/badge.svg" /></a>
<a href="https://github.com/layer5io/meshery-kuma/blob/master/LICENSE" alt="LICENSE">
  <img src="(https://img.shields.io/github/license/layer5io/meshery-kuma.svg" /></a>
<a href="https://github.com/layer5io/meshery-kuma/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22" alt="GitHub issues by-label">
  <img src="https://img.shields.io/github/issues/layer5io/meshery-kuma/help%20wanted.svg" /></a>
<a href="https://layer5.io/meshery" alt="Website">
  <img src="https://img.shields.io/website/https/layer5.io/meshery.svg" /></a>
<a href="https://twitter.com/intent/follow?screen_name=layer5" alt="Twitter Follow">
  <img src="https://img.shields.io/twitter/follow/layer5.svg?label=Follow&style=social" /></a>
<a href="https://slack.layer5.io" alt="Slack">
  <img src="https://img.shields.io/badge/Slack-@layer5.svg?logo=slack"></a>
<a href="https://bestpractices.coreinfrastructure.org/projects/3564" alt="CLI Best Practices">
  <img src="https://bestpractices.coreinfrastructure.org/projects/3564/badge" /></a>
</p>

# Meshery Adapter for Kuma

<p style="clear:both;">
<h2><a href="https://layer5.io/meshery">Meshery</a></h2>
<a href="https://meshery.io"><img src="img/readme/meshery-logo-light-text.svg"
style="margin:10px;" width="125px" 
alt="Meshery - the Service Mesh Management Plane" align="left" /></a>
<a href="https://meshery.io">Meshery</a> is the multi-service mesh management plane offering lifecycle management of more types of service meshes than any other tool available today. Meshery facilitates adopting, configuring, operating and managing performance of different service meshes and incorporates the collection and display of metrics from applications running on top of any service mesh. 
<br /><br /><p align="center"><i>If you’re using Meshery or if you like the project, please <a href="https://github.com/layer5io/meshery/stargazers">★</a> star this repository to show your support! 🤩</i></p>
</p>

<br />

<p style="clear:both;">
<h2><a name="contributing"></a><a name="community"></a> <a href="http://slack.layer5.io">Community</a> and <a href="https://github.com/layer5io/layer5/blob/master/CONTRIBUTING.md">Contributing</a></h2>
Our projects are community-built and welcome collaboration. 👍 Be sure to see the <a href="https://docs.google.com/document/d/17OPtDE_rdnPQxmk2Kauhm3GwXF1R5dZ3Cj8qZLKdo5E/edit">Layer5 Community Welcome Guide</a> for a tour of resources available to you and jump into our <a href="http://slack.layer5.io">Slack</a>! Contributors are expected to adhere to the <a href="https://github.com/cncf/foundation/blob/master/code-of-conduct.md">CNCF Code of Conduct</a>.


<a href="https://slack.meshery.io">

<picture align="right">
  <source media="(prefers-color-scheme: dark)" srcset="img\readme\slack-dark-128.png"  width="110px" align="right" style="margin-left:10px;margin-top:10px;">
  <source media="(prefers-color-scheme: light)" srcset="img\readme\slack-128.png" width="110px" align="right" style="margin-left:10px;padding-top:5px;">
  <img alt="Shows an illustrated light mode meshery logo in light color mode and a dark mode meshery logo dark color mode." src="img\readme\slack-128.png" width="110px" align="right" style="margin-left:10px;padding-top:13px;">
</picture>
</a>

<a href="https://meshery.io/community"><img alt="Layer5 Cloud Native Community" src="img/readme/community.svg" style="margin-right:8px;padding-top:5px;" width="140px" align="left" /></a>

<p>
✔️ <em><strong>Join</strong></em> <a href="https://drive.google.com/open?id=1c07UO9dS7_tFD-ClCWHIrEzRnzUJoFQ10EzfJTpS7FY">weekly community meeting</a> on <a href="https://bit.ly/2SbrRhe">Fridays from 10am - 11am Central</a>.<br />
✔️ <em><strong>Watch</strong></em> community <a href="https://www.youtube.com/channel/UCFL1af7_wdnhHXL1InzaMvA?sub_confirmation=1">meeting recordings</a>.<br />
✔️ <em>Fill-in</em> a <a href="https://layer5.io/newcomers">community member form</a> to gain access to community resources.<br />
✔️ <em><strong>Discuss</strong></em> in the <a href="https://discuss.layer5.io">Community Forum</a>.<br />
</p>
<p align="center">
<i>Not sure where to start?</i> Grab an open issue with the <a href="https://github.com/issues?q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Alayer5io+org%3Ameshery+org%3Aservice-mesh-performance+org%3Aservice-mesh-patterns+label%3A%22help+wanted%22+">help-wanted label</a>.
</p>

<br />
<br />

## About Layer5

**Community First**

<p>The <a href="https://layer5.io">Layer5</a> community represents the largest collection of service mesh projects and their maintainers in the world.</p>

**Open Source First**

<p>We build projects to provide learning environments, deployment and operational best practices, performance benchmarks, create documentation, share networking opportunities, and more. Our shared commitment to the open source spirit pushes Layer5 projects forward.</p>

**License**

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
