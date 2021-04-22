var app = Application.currentApplication()
app.includeStandardAdditions = true
app.activate()

ObjC.import('stdlib')
ObjC.import('readline')

function run(args) {
    Progress.totalUnitCount = 100
    Progress.completedUnitCount = 0
    Progress.description = args[0] || "Progress"
    Progress.additionalDescription = args[1] || "Running..."

    while (true) {
        var s
        try {
            s = $.readline('')
        } catch (e) {
            if (e.errorNumber === -128) $.exit(1)
            break
        }

        if (s.indexOf('#') === 0) {
            Progress.additionalDescription = s.slice(1).trim()
            continue
        }

        var i = parseInt(s)
        if (Number.isSafeInteger(i)) {
            Progress.completedUnitCount = i
            continue
        }
    }

    Progress.completedUnitCount = 100
}
