import * as React from 'react'
import * as ReactDOM from 'react-dom'
import './bootstrap.tsx'
import CssBaseline from '@material-ui/core/CssBaseline'
import Typography from '@material-ui/core/Typography'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'

import MAppBar from './components/MAppBar'


const theme = createMuiTheme({
  palette: {
    type: 'dark',
  },
})

const Main = ():JSX.Element => (
  <div id='root'>
    <MuiThemeProvider theme={theme}>
      <CssBaseline />
      <Typography variant='h2'>
        get it done
      </Typography>
    </MuiThemeProvider>
  </div>
)

ReactDOM.render(<Main />, document.getElementById('root'));
