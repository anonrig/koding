kd = require 'kd'
Encoder = require 'htmlencode'
MarkdownEditorView = require './markdowneditorview'
defaults = require 'app/util/stacks/defaults'
isMine  = require 'app/util/isMine'
isAdmin = require 'app/util/isAdmin'


module.exports = class ReadmeView extends kd.View

  constructor: (options = {}, data) ->

    super options, data

    { stackTemplate } = @getData()

    content = if stackTemplate?.description \
      then Encoder.htmlDecode stackTemplate?.description
      else defaults.description

    @editorView   = @addSubView new MarkdownEditorView
      content     : content
      delegate    : this
      contentType : 'md'


  viewAppended: ->

    super
    { stackTemplate } = @getData()
    isMine = isAdmin() or isMine(stackTemplate)
    @editorView.ready =>
      @setReadOnly()  unless isMine
      @listenEditorEvents()


  listenEditorEvents: ->
    @on 'FocusToEditor', => @editorView.setFocus yes


  setReadOnly: ->

    @setClass 'isntMine'
    @editorView.aceView.ace.editor.setReadOnly yes
