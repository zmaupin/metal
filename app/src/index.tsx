import * as React from 'react'
import * as ReactDOM from 'react-dom'
import CssBaseline from '@material-ui/core/CssBaseline'
import { MuiThemeProvider, createMuiTheme} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import MAppBar from './components/MAppBar'
import './index.tsx'

const theme = createMuiTheme({
  palette: {
    type: 'dark',
  },
})

const Main = ():JSX.Element => (
  <div id='root'>
    <MuiThemeProvider theme={theme}>
      <CssBaseline />
      <MAppBar title='Rexec'/>
      <Typography variant='h2'>
        get it done
      </Typography>
    </MuiThemeProvider>
  </div>
)

ReactDOM.render(<Main />, document.getElementById('root'));
