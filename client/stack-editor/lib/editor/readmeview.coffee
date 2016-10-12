kd = require 'kd'
Encoder = require 'htmlencode'
MarkdownEditorView = require './markdowneditorview'
StackBaseEditorTabView = require './stackbaseeditortabview'


module.exports = class ReadmeView extends StackBaseEditorTabView

  constructor: (options = {}, data) ->

    super options, data

    { stackTemplate } = @getData()

    content = if stackTemplate?.description \
      then Encoder.htmlDecode stackTemplate?.description
      else ''

    @editorView   = @addSubView new MarkdownEditorView
      content     : content
      delegate    : this
      contentType : 'md'
