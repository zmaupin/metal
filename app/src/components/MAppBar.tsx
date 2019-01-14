import * as React from 'react'
import AppBar from '@material-ui/core/AppBar'
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton'
import MenuIcon from '@material-ui/icons/Menu'
import Toolbar from '@material-ui/core/Toolbar'
import Typography from '@material-ui/core/Typography'
import { createStyles, withStyles, WithStyles } from '@material-ui/styles'
import { Theme } from '@material-ui/core/styles/createMuiTheme'


const styles = (theme: Theme) => createStyles({
  root: {
    flexGrow: 1,
  },
  grow: {
    flexGrow: 1,
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 20,
  },
})

export interface IMAppBarProps extends WithStyles<typeof styles> {
  title: string
  color?: string
}

const MAppBar = withStyles(styles)((props: IMAppBarProps) => {
  return (
    <div className={props.classes.root}>
      <AppBar position="static" style={{backgroundColor: props.color ? props.color : '#5A7E8A'}}>
        <Toolbar>
          <IconButton className={props.classes.menuButton} color="inherit" aria-label="Menu">
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" color="inherit" className={props.classes.grow}>
            {props.title}
          </Typography>
          <Button color="inherit">Metal</Button>
        </Toolbar>
      </AppBar>
    </div>
  )
})

export default MAppBar
