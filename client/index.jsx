
import React from 'react'
import ReactDOM from 'react-dom'

class AjaxAssistant {

  constructor($) {
    this.$ = $
  }

  handleError(xhr, reject) {
    var text = ""
    try {
      const data = JSON.parse(xhr.responseText)
      text = data.errors
    } catch(e) {
      text = xhr.responseText
    }

    if(text === "") {
      if(xhr.status === 404) {
        text = "A resource could not be found."
      }
    }
    
    reject(text)
  }
  
  post(path, data) {
    return new Promise((resolve, reject) => {
      if(!data) data = {}
      this.$.ajax({
        method: 'POST',
        url: path,
        dataType: 'JSON',
        data: data,
        beforeSend: (xhr) => { xhr.setRequestHeader('CSRF-Token', this.$('meta[name=csrf-token]').attr('content')) },
        success: (data) => resolve(data),
        error: (xhr) => this.handleError(xhr, reject)
      })
    })
  }
  
  get(path) {
    return new Promise((resolve, reject) => {
      this.$.ajax({
        url: path,
        dataType: 'JSON',
        success: (data) => resolve(data),
        error: (xhr) => this.handleError(xhr, reject)
      })
    })
  }
}

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      gists: null
    }

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
    const gists = this.state.gists ? this.state.gists.map((g) => { return (
      <tr key={g.id}>
        <td>{g.title}</td><td><a href={this.gistUrl(g)} target="_blank">{g.vendor_id}</a></td>
      </tr>
    )}) : ""

    const login = this.props.loggedIn ?
          (<span><a href="/api/gists/fetchAll">Fetch Gists</a> | {this.props.username} <a href="/logout">Logout</a></span>)
          : (<a href="/oauth/github">Login with Github</a>)

    return (
      <div className="gists">
        <div className="github-login">
          {login}
        </div>
        <table className="table table-bordered">
          <thead><tr><td>Name of Gist</td><td>Link</td></tr></thead>
          <tbody>{gists}</tbody>
        </table>
      </div>
    )
  }
}

$(() => {
  ReactDOM.render(<App username={username} loggedIn={loggedIn}/>, document.getElementById('app-root'))
})
