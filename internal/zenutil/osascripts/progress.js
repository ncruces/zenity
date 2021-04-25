var app = Application.currentApplication()
app.includeStandardAdditions = true
app.activate()

ObjC.import('stdlib')
ObjC.import('readline')

try { Progress.totalUnitCount = $.getenv('total') } catch { }
try { Progress.description = $.getenv('description') } catch { }

while (true) {
    var s
    try {
        s = $.readline('')
    } catch (e) {
        if (e.errorNumber === -128) $.exit(1)
        break
    }

    if (s.indexOf('#') === 0) {
        Progress.additionalDescription = s.slice(1)
        continue
    }

    var i = parseInt(s)
    if (i >= 0 && Progress.totalUnitCount > 0) {
        Progress.completedUnitCount = i
        continue
    }
}