class App extends React.Component {
  constructor(props) {
    super(props)
  }
  
  render() {
    const gists = this.props.gists ? this.props.gists.map(function(g) { return (
      <tr key={g.id}>
        <th>{g.title}</th><th><a href={g.url}>{g.url}</a></th>
      </tr>
    )}) : ""

    const login = this.props.loggedIn ?
          (<span><a href="/api/gists/fetchAll">Fetch Gists</a> | <a href="/logout">Logout</a></span>)
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
  ReactDOM.render(<App loggedIn={loggedIn}/>, document.getElementById('app-root'))
})
