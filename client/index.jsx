
import React from 'react'
import ReactDOM from 'react-dom'
import AjaxAssistant from 'AjaxAssistant.jsx'

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      query: "",
      gists: null,
      results: null,
      busy: false
    }

    this.onSubmit = this.onSubmit.bind(this)
    this.onChange = this.onChange.bind(this)
  }

  onChange(e) {
    const target = e.target
    const name = target.name
    this.setState({[name]: target.value})
  }
  
  onSubmit(e) {
    e.preventDefault()

    if(this.state.busy) return
    this.setState({busy: true})
    
    new AjaxAssistant($).get('/api/gists/search', {q: this.state.query})
      .then(results => this.setState({busy: false, results: results}))
  }
  
  componentDidMount() {
    if(!this.state.gists && this.props.loggedIn) {
      new AjaxAssistant($).get('/api/gists')
        .then(gists => this.setState({gists: gists}))
    }
  }

  gistUrl(gist) {
    return `https://gist.github.com/${this.props.username}/${gist.vendor_id}`
  }
  
  render() {
    const login = this.props.loggedIn ?
          (<span><a href="/api/gists/fetchAll">Fetch Gists</a> | {this.props.username} <a href="/logout">Logout</a></span>)
          : (<a href="/oauth/github">Login with Github</a>)

    let coreContent = ""
    if(this.state.results === null) {
      const gists = this.state.gists ? this.state.gists.map((g) => { return (
        <tr key={g.id}>
          <td>{g.title}</td>
          <td>{g.filename}</td>
          <td><a href={this.gistUrl(g)} target="_blank">{g.vendor_id}</a></td>
        </tr>
      )}) : <tr><td colSpan="3"></td></tr>

      const listAllGists = (
        <table className="table table-bordered">
          <thead><tr><th>Gist Description</th><th>Filename</th><th>ID / Link</th></tr></thead>
          <tbody>
            {gists}
          </tbody>
        </table>      
      )

      coreContent = listAllGists
    } else {
      const snippets = this.state.results.map((snippet, i) => { return (
        <tr key={i}>
          <td>
            <code>{snippet.body}</code>
          </td>
          <td></td>
          <td><a href={this.gistUrl(snippet)} target="_blank">{snippet.vendor_id}</a></td>
        </tr>
      )})
            
      const listResults = (
        <table className="table table-bordered">
          <thead><tr><th>Snippet</th><th>ID / Link</th></tr></thead>
          <tbody>
            {snippets}
          </tbody>
        </table>      
      )
      
      coreContent = listResults
    }
    
    return (
      <div className="gists">
        <div className="github-login">
          {login}
        </div>

        <form onSubmit={this.onSubmit}>
          <input type="text" onChange={this.onChange} value={this.state.query} name="query" placeholder="Search"/>
        </form>

        {coreContent}
      </div>
    )
  }
}

$(() => {
  ReactDOM.render(<App username={username} loggedIn={loggedIn}/>, document.getElementById('app-root'))
})
