class App extends React.Component {
  constructor(props) {
    super(props)
  }
  
  render() {
    const gists = this.props.gists.map(function(g) { return (
      <tr key={g.id}>
        <th>{g.title}</th><th><a href={g.url}>{g.url}</a></th>
      </tr>
    )})
    
    return (
      <div className="gists">
        <div className="github-login">
          <a href="/oauth/github">Login with Github</a>
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
  ReactDOM.render(<App gists={__bootstrap}/>, document.getElementById('app-root'))
})
