class TerminalModal extends KDModalView

  constructor: (options={}, data)->

    options.cssClass  = KD.utils.curry "terminal", options.cssClass

    if options.machine?.getName? and not options.title
      options.title = "Terminal on #{options.machine.getName()}"

    super options, data

    {machine, command} = @getOptions()

    @webterm           = new WebTermView
      delegate         : this
      cssClass         : "webterm"
      machine          : machine
      advancedSettings : no

    @webterm.on "WebTermEvent", (data)=>
      @emit "terminal.event", data

    @webterm.once "WebTermConnected", (remote)=>

      @emit "terminal.connected", remote

      @on "terminal.input", (command)=>
        remote.input command
        @webterm.getDomElement().click()

      @run command  if command?

    @webterm.on "WebTerm.terminated", =>
      @emit "terminal.terminated"

    @on "click", => @setCss 'zIndex', 10000 + KD.utils.uniqueId()

  viewAppended: ->

    @addSubView @webterm
    @webterm.connectToTerminal()
    @setCss 'zIndex', 10000 + KD.utils.uniqueId()

  run: (command)->

    @emit "terminal.input", command + "\n"
